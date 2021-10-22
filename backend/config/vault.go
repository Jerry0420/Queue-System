package config

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/jerry0420/queue-system/backend/logging"
)

type vaultWrapper struct {
	credName string
	client *api.Client
	logger logging.LoggerTool
}

func NewVault(address string, roleID string, wrappedToken string, credName string, logger logging.LoggerTool) *vaultWrapper {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)
	if err != nil {
		logger.FATALf("Fail to create connection with vault server.")
	}

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

	return &vaultWrapper{credName, client, logger}
}

func (vault *vaultWrapper) checkAndRenewToken() {
	var tokenInfo *api.Secret
	var err error
	var ttl time.Duration
	token := vault.client.Auth().Token()

	for {
		tokenInfo, err = token.LookupSelf()
		if err != nil {
			vault.logger.ERRORf("Fail to lookup token info. %v", err)
		}
		ttl, err = tokenInfo.TokenTTL()
		if err != nil {
			vault.logger.ERRORf("Fail to get token ttl. %v", err)
		}
		if ttl <= time.Minute * 30 {
			tokenInfo, err = token.RenewSelf(3600)
			if err != nil {
				vault.logger.ERRORf("Fail to renew token. %v", err)
			}
		} else {
			// May be some delay after the server running long period of time.
			time.Sleep(ttl - time.Minute * 30)
		}
	}
}

func (vault *vaultWrapper) checkAndRenewCred(leaseId string) {
	var credInfo *api.Secret
	var err error
	var ttl time.Duration
	sys := vault.client.Sys()

	for {
		credInfo, err = sys.Lookup(leaseId)
		if err != nil {
			vault.logger.ERRORf("Fail to lookup cred info. %v", err)
		}
		ttl, err = credInfo.TokenTTL()
		if err != nil {
			vault.logger.ERRORf("Fail to get cred ttl. %v", err)
		}
		if ttl <= time.Minute * 30 {
			credInfo, err = sys.Renew(leaseId, 3600)
			if err != nil {
				vault.logger.ERRORf("Fail to renew cred %s %v", leaseId, err)
			}
		} else {
			// May be some delay after the server running long period of time.
			time.Sleep(ttl - time.Minute * 30)
		}
	}
}

func (vault *vaultWrapper) GetDbCred() (string, string, string) {
	credPath := fmt.Sprintf("database/creds/%s", vault.credName)
	cred, err := vault.client.Logical().Read(credPath)
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

	leaseID := cred.LeaseID

	go vault.checkAndRenewToken()
	go vault.checkAndRenewCred(leaseID)

	return username, password, leaseID
}

func (vault *vaultWrapper) RevokeLeaseAndToken(leaseId string) {
	err := vault.client.Sys().Revoke(leaseId)
	if err != nil {
		vault.logger.WARNf("Fail to revoke lease id %s %v", leaseId, err)
	}
	err = vault.client.Auth().Token().RevokeSelf("whatever it is.")
	if err != nil {
		vault.logger.WARNf("Fail to revoke token. %v", err)
	}
}