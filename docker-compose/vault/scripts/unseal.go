package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// This script will unseal vault server and get a pair of role id and wrapped token.

// manually run this script when...
// 1. vault server restart.
// 2. vault server start.

// After get a pair of role id and wrapped token, copy and paste to server env file, the server will then unwrap the wrapped token to get the secret id. 
// With role id and secret id, the server can login vault server and retrive a periodic token to interact with vault server. 

func unseal() {
	const basePath = "/vault"
	configPath := filepath.Join(basePath, "config")
	policiesPath := filepath.Join(basePath, "policies")
	logsPath := filepath.Join(basePath, "logs")

	dbPolicyFile := filepath.Join(policiesPath, "db_hcl.sh")
	dbConnectionFile := filepath.Join(configPath, "connection_db.sh")
	ruleSqlFile := filepath.Join(configPath, "create_user.sql")
	logFile := filepath.Join(logsPath, "audit.log")

	dbName := os.Getenv("POSTGRES_BACKEND_DB")
	roleName := os.Getenv("VAULT_ROLE_NAME")
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

	cmd = fmt.Sprintf(
		"vault login %s", 
		rootToken,
	)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		panic("Fail to login with root token.")
	}

	cmd = "vault secrets enable database"
	exec.Command("sh", "-c", cmd).Output()

	cmd = fmt.Sprintf(
		"vault audit enable file file_path=%s", 
		logFile,
	)
	exec.Command("sh", "-c", cmd).Output()

	cmd = "vault auth enable approle"
	exec.Command("sh", "-c", cmd).Output()

	// ============================================

	cmd = fmt.Sprintf(
		"%s | vault policy write %s -", 
		dbPolicyFile, 
		policyName,
	)
	exec.Command("sh", "-c", cmd).Output()

	cmd = fmt.Sprintf(
		"vault write auth/approle/role/%s token_policies=\"%s\" token_ttl=%s token_max_ttl=%s token_num_uses=%d secret_id_ttl=%s secret_id_num_uses=%d", 
		roleName, 
		policyName,
		"1h", // token_ttl 
		"24h", // token_max_ttl
		0, // token_num_uses
		"30m", // secret_id_ttl
		1, // secret_id_num_uses
	)
	exec.Command("sh", "-c", cmd).Output()

	// ===========================================

	cmd = fmt.Sprintf(
		"%s | vault write database/config/%s -", 
		dbConnectionFile, 
		dbName,
	)
	exec.Command("sh", "-c", cmd).Output()

	// cmd = fmt.Sprintf("vault write -force database/rotate-root/%s", dbName)
	// exec.Command("sh", "-c", cmd).Output()

	cmd = fmt.Sprintf(
		"vault write database/roles/%s db_name=%s creation_statements=@%s default_ttl=%s max_ttl=%s", 
		roleName, 
		dbName, 
		ruleSqlFile, 
		"1h", // default_ttl of db cred
		"24h", // max_ttl of db cred
	)
	exec.Command("sh", "-c", cmd).Output()

	// =============================================

	cmd = fmt.Sprintf(
		"vault read auth/approle/role/%s/role-id -format=json", 
		roleName,
	)
	out, _ = exec.Command("sh", "-c", cmd).Output()
	var appRoleResults map[string]interface{}
	json.Unmarshal(out, &appRoleResults)

	fmt.Println("\n=========================Role id====================================")
	fmt.Println(appRoleResults["data"].(map[string]interface{})["role_id"])
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