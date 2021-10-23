package config

import (
	"fmt"
	"time"
	"net/http"
	"bytes"
    "encoding/json"

	"github.com/hashicorp/vault/api"
	"github.com/jerry0420/queue-system/backend/logging"
)

func NewVaultConnection(vaultAddress string, vaultWrappedTokenAddress string, roleID string, credName string, logger logging.LoggerTool) (*api.Logical, *api.TokenAuth, *api.Sys) {
	config := api.DefaultConfig()
	config.Address = vaultAddress

	client, err := api.NewClient(config)
	if err != nil {
		logger.FATALf("Fail to create connection with vault server.")
	}

	params := map[string]string{"roleName": credName}
	json_params, _ := json.Marshal(params)
	httpClient := http.Client{Timeout: 10 * time.Second}
	
	// Everytime, when server start up, get wrapped token from vault server.
	response, err := httpClient.Post(vaultWrappedTokenAddress + "/wrappedToken", "application/json", bytes.NewBuffer(json_params))
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
		"role_id":   roleID,
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

type vaultWrapper struct {
	logical *api.Logical
	token *api.TokenAuth
	sys *api.Sys
	leaseID string
	logger logging.LoggerTool
}

func NewVaultWrapper(logical *api.Logical, token *api.TokenAuth, sys *api.Sys, logger logging.LoggerTool) *vaultWrapper {
	return &vaultWrapper{
		logical: logical,
		token: token,
		sys: sys,
		logger: logger,
	}
}

func (vault *vaultWrapper) checkAndRenewToken() {
	var tokenInfo *api.Secret
	var err error
	var ttl time.Duration

	for {
		tokenInfo, err = vault.token.LookupSelf()
		if err != nil {
			vault.logger.ERRORf("Fail to lookup token info. %v", err)
		}
		ttl, err = tokenInfo.TokenTTL()
		if err != nil {
			vault.logger.ERRORf("Fail to get token ttl. %v", err)
		}
		if ttl <= time.Minute * 30 {
			tokenInfo, err = vault.token.RenewSelf(3600)
			if err != nil {
				vault.logger.ERRORf("Fail to renew token. %v", err)
			}
		} else {
			// May be some delay after the server running long period of time.
			time.Sleep(ttl - time.Minute * 30)
		}
	}
}

func (vault *vaultWrapper) checkAndRenewCred() {
	var credInfo *api.Secret
	var err error
	var ttl time.Duration

	for {
		credInfo, err = vault.sys.Lookup(vault.leaseID)
		if err != nil {
			vault.logger.ERRORf("Fail to lookup cred info. %v", err)
		}
		ttl, err = credInfo.TokenTTL()
		if err != nil {
			vault.logger.ERRORf("Fail to get cred ttl. %v", err)
		}
		if ttl <= time.Minute * 30 {
			credInfo, err = vault.sys.Renew(vault.leaseID, 3600)
			if err != nil {
				vault.logger.ERRORf("Fail to renew cred %s %v", vault.leaseID, err)
			}
		} else {
			// May be some delay after the server running long period of time.
			time.Sleep(ttl - time.Minute * 30)
		}
	}
}

func (vault *vaultWrapper) GetDbCred(credName string) (string, string) {
	credPath := fmt.Sprintf("database/creds/%s", credName)
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

	vault.leaseID = cred.LeaseID

	go vault.checkAndRenewToken()
	go vault.checkAndRenewCred()

	return username, password
}

func (vault *vaultWrapper) RevokeLeaseAndToken() {
	err := vault.sys.Revoke(vault.leaseID)
	if err != nil {
		vault.logger.WARNf("Fail to revoke lease id %s %v", vault.leaseID, err)
	}
	err = vault.token.RevokeSelf("whatever it is.")
	if err != nil {
		vault.logger.WARNf("Fail to revoke token. %v", err)
	}
}