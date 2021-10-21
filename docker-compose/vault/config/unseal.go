package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// This script will unseal vault server and get a periodic token.

// manually run this script when...
// 1. vault server restart.
// 2. vault server start.

// After get a periodic token, copy and paste to server env file.

func unseal() {
	const basePath = "/vault"
	configPath := filepath.Join(basePath, "config")
	policiesPath := filepath.Join(basePath, "policies")

	dbPolicyFile := filepath.Join(policiesPath, "db_hcl.sh")
	dbConnectionFile := filepath.Join(configPath, "connection_db.sh")
	ruleSqlFile := filepath.Join(configPath, "create_user.sql")

	dbName := os.Getenv("POSTGRES_BACKEND_DB")
	credName := os.Getenv("VAULT_CRED_NAME")
	policyName := os.Getenv("VAULT_POLICY_NAME")

	cmd := "vault operator init -format=json -status=true"
	out, _ := exec.Command("sh", "-c", cmd).Output()
	var unsealStatus map[string]interface{}
	json.Unmarshal(out, &unsealStatus)

	var unsealKeys []interface{}
	var rootToken string

	if unsealStatus["Initialized"].(bool) == false {
		cmd = "vault operator init -format=json"
		out, _ = exec.Command("sh", "-c", cmd).Output()
		var initResults map[string]interface{}
		json.Unmarshal(out, &initResults)
		unsealKeys = initResults["unseal_keys_b64"].([]interface{})
		rootToken = initResults["root_token"].(string)
		fmt.Println("=========================Unseal Keys====================================")
		fmt.Println(unsealKeys)
		fmt.Println()
		
		fmt.Println("=========================Root Token====================================")
		fmt.Println(rootToken)
		fmt.Println()
		
	} else {
		fmt.Println("Please enter unseal keys in a slice.")
		var input string
    	fmt.Scanln(&input)
		var initResults []interface{}
		json.Unmarshal([]byte(input), &initResults)
		unsealKeys = initResults
		
		fmt.Println("Please enter root token.")
		fmt.Scanln(&rootToken)
	}

	for index, unsealKey := range unsealKeys {
		if index >= 3 {
			break
		}
		cmd = fmt.Sprintf("vault operator unseal %s", unsealKey)
		_, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			panic("Fail to unseal vault server.")
		}
	}

	cmd = fmt.Sprintf("vault login %s", rootToken)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		panic("Fail to login with root token.")
	}

	cmd = "vault secrets enable database"
	exec.Command("sh", "-c", cmd).Output()

	cmd = fmt.Sprintf("%s | vault policy write %s -", dbPolicyFile, policyName)
	exec.Command("sh", "-c", cmd).Output()

	cmd = fmt.Sprintf("%s | vault write database/config/%s -", dbConnectionFile, dbName)
	exec.Command("sh", "-c", cmd).Output()

	// cmd = fmt.Sprintf("vault write -force database/rotate-root/%s", dbName)
	// exec.Command("sh", "-c", cmd).Output()

	cmd = fmt.Sprintf("vault write database/roles/%s db_name=%s creation_statements=@%s default_ttl=1h max_ttl=24h", credName, dbName, ruleSqlFile)
	exec.Command("sh", "-c", cmd).Output()

	cmd = fmt.Sprintf("vault token create -policy=%s -period 1h -format=json", policyName)
	out, _ = exec.Command("sh", "-c", cmd).Output()
	var dbTokenResults map[string]interface{}
	json.Unmarshal(out, &dbTokenResults)

	fmt.Println("\n=========================DB Token====================================")
	fmt.Println(dbTokenResults["auth"].(map[string]interface{})["client_token"])
	fmt.Println()
}

func doUnseal(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			syscall.Kill(1, syscall.SIGTERM)
		}
	}()
	fn()
}

func main() {
	doUnseal(unseal)
}