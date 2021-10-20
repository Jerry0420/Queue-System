#!/usr/bin/env python

import json
import subprocess
import os

'''
This script will unseal vault server and get a periodic token.

manually run this script when...
1. vault server restart.
2. vault server start.

After get a periodic token, copy and paste to server env file.
'''

base_path = "/vault"
file_path = os.path.join(base_path, "file")
config_path = os.path.join(base_path, "config")
policies_path = os.path.join(base_path, "policies")

unseal_key_file = os.path.join(file_path, 'keys.json')
db_policy_file = os.path.join(policies_path, "db_hcl.sh")
db_connection_file = os.path.join(config_path, "connection_db.sh")
rule_sql_file = os.path.join(config_path, "create_user.sql")
db_token_file = os.path.join(file_path, "db_token.json")

db_name = os.getenv("POSTGRES_BACKEND_DB")
cred_name = os.getenv("VAULT_CRED_NAME")
policy_name = os.getenv("VAULT_POLICY_NAME")

if __name__ == '__main__':
    if not os.path.isfile(unseal_key_file):
        subprocess.call(f'vault operator init -format=json > {unseal_key_file}', shell=True)

    data = {}
    root_token = None
    with open(unseal_key_file, 'r') as openfile: 
        data = json.load(openfile) 
    unseal_keys = data['unseal_keys_b64']
    root_token = data['root_token']

    for key, unseal_key in enumerate(unseal_keys):
        if key >= 3:
            break
        subprocess.call(f'vault operator unseal {unseal_key}', shell=True)

    subprocess.call(f'vault login {root_token}', shell=True)
    subprocess.call('vault secrets enable database', shell=True)

    subprocess.call(f'{db_policy_file} | vault policy write {policy_name} -', shell=True)
    subprocess.call(f'{db_connection_file} | vault write database/config/{db_name} -', shell=True)
    subprocess.call(f'vault write database/roles/{cred_name} db_name={db_name} creation_statements=@{rule_sql_file} default_ttl=1h max_ttl=24h', shell=True)
    
    subprocess.call(f'vault token create -policy={policy_name} -period 1h -format=json > {db_token_file}', shell=True)

    if os.path.isfile(db_token_file):
      with open(db_token_file, 'r') as openfile: 
        data = json.load(openfile) 
        print('======================db_token==============================')
        print(data['auth']['client_token'])
        print('======================db_token==============================')
        os.remove(db_token_file)
    else:
        print("Failed, please run again!")