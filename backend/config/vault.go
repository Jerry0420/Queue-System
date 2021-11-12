package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/jerry0420/queue-system/backend/logging"
)

type VaultConnectionConfig struct {
	Address string
	WrappedTokenAddress string
	RoleID string
	CredName string
}

func NewVaultConnection(logger logging.LoggerTool, vaultConnectionConfig *VaultConnectionConfig) (*api.Logical, *api.TokenAuth, *api.Sys) {
	config := api.DefaultConfig()
	config.MaxRetries = 5
	config.Address = vaultConnectionConfig.Address

	client, err := api.NewClient(config)
	if err != nil {
		logger.FATALf("Fail to create connection with vault server.")
	}

	params := map[string]string{"roleName": vaultConnectionConfig.CredName}
	json_params, _ := json.Marshal(params)
	httpClient := http.Client{Timeout: 10 * time.Second}
	
	// Everytime, when server start up, get wrapped token from vault server.
	response, err := httpClient.Post(vaultConnectionConfig.WrappedTokenAddress + "/wrappedToken", "application/json", bytes.NewBuffer(json_params))
	if err != nil || response.StatusCode != http.StatusOK{
        logger.FATALf("Fail to get wrapped token.")
    }

	var decodedResponse map[string]interface{}
    json.NewDecoder(response.Body).Decode(&decodedResponse)
	wrappedToken := decodedResponse["wrappedToken"].(string)

	logical := client.Logical()
	unwrappedToken, err := logical.Unwrap(wrappedToken)
	if err != nil {
		logger.FATALf("Fail to unwrap token. %v", err)
	}
	secretID := unwrappedToken.Data["secret_id"]

	vaulrRoleLoginParams := map[string]interface{}{
		"role_id":   vaultConnectionConfig.RoleID,
		"secret_id": secretID,
	}
	loginResponse, err := logical.Write("auth/approle/login", vaulrRoleLoginParams)
	if err != nil {
		logger.FATALf("Fail to login with approle. %v", err)
	}
	if loginResponse == nil || loginResponse.Auth == nil || loginResponse.Auth.ClientToken == "" {
		logger.FATALf("Fail to retrive login info. %v", err)
	}

	client.SetToken(loginResponse.Auth.ClientToken)

	token := client.Auth().Token()
	sys := client.Sys()
	
	return logical, token, sys
}

type VaultWrapper struct {
	logical *api.Logical
	token *api.TokenAuth
	sys *api.Sys
	credName string
	logger logging.LoggerTool
}

func NewVaultWrapper(credName string, logical *api.Logical, token *api.TokenAuth, sys *api.Sys, logger logging.LoggerTool) *VaultWrapper {
	vault := &VaultWrapper{
		credName: credName,
		logical: logical,
		token: token,
		sys: sys,
		logger: logger,
	}
	
	go vault.checkAndRenewToken()

	return vault 
}

func (vault *VaultWrapper) checkAndRenewToken() {
	var tokenInfo *api.Secret
	var err error
	var ttl time.Duration

	for {
		tokenInfo, err = vault.token.LookupSelf()
		if err != nil {
			vault.logger.FATALf("Fail to lookup token info. %v", err)
		}
		ttl, _ = tokenInfo.TokenTTL()
		if ttl <= time.Minute * 30 {
			_, err = vault.token.RenewSelf(3600)
			if err != nil {
				vault.logger.FATALf("Fail to renew token. %v", err)
			}
		} else {
			// May be some delay after the server running long period of time.
			time.Sleep(ttl - time.Minute * 30)
		}
	}
}

func (vault *VaultWrapper) checkAndRenewCred(leaseID string, credExpireChan chan bool, leaseRevocableChan chan bool) {
	var credInfo *api.Secret
	var currentExpireTime time.Time
	var renewExpireTime time.Time
	var err error
	var ttl time.Duration

	for {
		credInfo, err = vault.sys.Lookup(leaseID)
		if err != nil {
			vault.logger.FATALf("Fail to lookup cred info. %v", err)
		}
		currentExpireTime, _ = time.Parse(time.RFC3339Nano, credInfo.Data["expire_time"].(string))
		ttl, _ = credInfo.TokenTTL()
		
		if ttl <= time.Minute * 30 {
			_, err = vault.sys.Renew(leaseID, 3600)
			if err != nil {
				vault.logger.FATALf("Fail to renew cred %s %v", leaseID, err)
			}
			credInfo, err = vault.sys.Lookup(leaseID)
			if err != nil {
				vault.logger.FATALf("Fail to lookup cred info. %v", err)
			}
			renewExpireTime, _ = time.Parse(time.RFC3339Nano, credInfo.Data["expire_time"].(string))
			// when renewExpireTime and currentExpireTime are approximately the same, mark this cred as expired!
			if renewExpireTime.Sub(currentExpireTime) <= time.Minute * 1 {
				credExpireChan <- true
				<- leaseRevocableChan
				vault.revokeLease(leaseID)
				return
			}
		} else {
			// May be some delay after the server running long period of time.
			time.Sleep(ttl - time.Minute * 30)
		}
	}
}

func (vault *VaultWrapper) GetDbCred(leaseRevocableChan chan bool) (string, string, chan bool) {
	credExpireChan := make(chan bool, 1)
	credPath := fmt.Sprintf("database/creds/%s", vault.credName)
	cred, err := vault.logical.Read(credPath)
	if err != nil {
		vault.logger.FATALf("Fail to get cred in vault. %v", err)
	}

	username, ok := cred.Data["username"].(string)
	if !ok {
		vault.logger.FATALf("Fail to get username in vault.")
	}

	password, ok := cred.Data["password"].(string)
	if !ok {
		vault.logger.FATALf("Fail to get password in vault.")
	}

	go vault.checkAndRenewCred(cred.LeaseID, credExpireChan, leaseRevocableChan)

	return username, password, credExpireChan
}

func (vault *VaultWrapper) revokeLease(leaseID string) {
	err := vault.sys.Revoke(leaseID)
	if err != nil {
		vault.logger.WARNf("Fail to revoke lease id %s %v", leaseID, err)
	}
}

func (vault *VaultWrapper) RevokeToken() error {
	err := vault.token.RevokeSelf("whatever it is.")
	return err
}