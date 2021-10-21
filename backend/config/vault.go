package config

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/jerry0420/queue-system/backend/logging"
)

type vaultWrapper struct {
	address string
	token string
	credName string
	client *api.Client
	logger logging.LoggerTool
}

func NewVault(address string, token string, credName string, logger logging.LoggerTool) *vaultWrapper {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)
	if err != nil {
		logger.FATALf("Fail to create connection with vault server.")
	}
	return &vaultWrapper{address, token, credName, client, logger}
}

func (vault *vaultWrapper) checkAndRenewToken() {
	for {
		tokenInfo, err := vault.client.Auth().Token().LookupSelf()
		if err != nil {
			vault.logger.ERRORf("Fail to lookup token info. %v", err)
		}
		ttl, err := tokenInfo.TokenTTL()
		if err != nil {
			vault.logger.ERRORf("Fail to get token ttl. %v", err)
		}
		if ttl <= time.Minute * 15 {
			tokenInfo, err = vault.client.Auth().Token().RenewSelf(3600)
			if err != nil {
				vault.logger.ERRORf("Fail to renew token. %v", err)
			}
		} else {
			timer := time.NewTimer(ttl - time.Minute * 15)
			<- timer.C
		}
	}
}

func (vault *vaultWrapper) checkAndRenewCred(leaseId string) {
	for {
		credInfo, err := vault.client.Sys().Lookup(leaseId)
		if err != nil {
			vault.logger.ERRORf("Fail to lookup cred info. %v", err)
		}
		ttl, err := credInfo.TokenTTL()
		if err != nil {
			vault.logger.ERRORf("Fail to get cred ttl. %v", err)
		}
		if ttl <= time.Minute * 15 {
			credInfo, err = vault.client.Sys().Renew(leaseId, 3600)
			if err != nil {
				vault.logger.ERRORf("Fail to renew cred %s %v", leaseId, err)
			}
		} else {
			timer := time.NewTimer(ttl - time.Minute * 15)
			<- timer.C
		}
	}
}

func (vault *vaultWrapper) GetDbSecret() (string, string, string) {
	vault.client.SetToken(vault.token)

	cred := fmt.Sprintf("database/creds/%s", vault.credName)
	secret, err := vault.client.Logical().Read(cred)
	if err != nil {
		vault.logger.FATALf("Fail to get cred in vault. %v", err)
	}

	username, ok := secret.Data["username"].(string)
	if !ok {
		vault.logger.FATALf("Fail to get username in vault.")
	}

	password, ok := secret.Data["password"].(string)
	if !ok {
		vault.logger.FATALf("Fail to get password in vault.")
	}

	leaseID := secret.LeaseID

	go vault.checkAndRenewToken()
	go vault.checkAndRenewCred(leaseID)

	return leaseID, username, password
}

func (vault *vaultWrapper) RevokeLease(leaseId string) {
	err := vault.client.Sys().Revoke(leaseId)
	if err != nil {
		vault.logger.WARNf("Fail to revoke lease id %s %v", leaseId, err)
	}
}