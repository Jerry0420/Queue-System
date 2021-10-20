#!/usr/bin/env python

import json
import subprocess
import os
import signal

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

db_policy_file = os.path.join(policies_path, "db_hcl.sh")
db_connection_file = os.path.join(config_path, "connection_db.sh")
rule_sql_file = os.path.join(config_path, "create_user.sql")

db_name = os.getenv("POSTGRES_BACKEND_DB")
cred_name = os.getenv("VAULT_CRED_NAME")
policy_name = os.getenv("VAULT_POLICY_NAME")

def main():
    pipe = subprocess.Popen('vault operator init -format=json -status=true', shell=True, stdout=subprocess.PIPE)
    unseal_status = pipe.communicate()[0]
    unseal_status = json.loads(unseal_status)
    unseal_status = unseal_status['Initialized']

    unseal_keys = []
    root_token = None

    if unseal_status == False:
        pipe = subprocess.Popen(f'vault operator init -format=json', shell=True, stdout=subprocess.PIPE)
        init_results = pipe.communicate()[0]
        init_results = json.loads(init_results)
        unseal_keys = init_results['unseal_keys_b64']
        root_token = init_results['root_token']
        print('================================unseal_keys========================================')
        print(unseal_keys)
        print('================================unseal_keys========================================\n\n')
        print('================================root_token========================================')
        print(root_token)        
        print('================================root_token========================================\n\n')
    else:
        init_results = input("Please input unseal keys in a list. : \n")
        init_results = json.loads(init_results.replace("'", '"'))
        unseal_keys = init_results
        root_token = input("Please input root token. : \n")

    for key, unseal_key in enumerate(unseal_keys):
        if key >= 3:
            break
        subprocess.call(f'vault operator unseal {unseal_key}', shell=True)

    subprocess.call(f'vault login {root_token}', shell=True)
    subprocess.call('vault secrets enable database', shell=True)

    subprocess.call(f'{db_policy_file} | vault policy write {policy_name} -', shell=True)
    subprocess.call(f'{db_connection_file} | vault write database/config/{db_name} -', shell=True)
    # subprocess.call(f'vault write -force database/rotate-root/{db_name}', shell=True)
    subprocess.call(f'vault write database/roles/{cred_name} db_name={db_name} creation_statements=@{rule_sql_file} default_ttl=1h max_ttl=24h', shell=True)
    
    pipe = subprocess.Popen(f'vault token create -policy={policy_name} -period 1h -format=json', shell=True, stdout=subprocess.PIPE)
    db_token_results = pipe.communicate()[0]
    db_token_results = json.loads(db_token_results)

    print('======================db_token==============================')
    print(db_token_results['auth']['client_token'])
    print('======================db_token==============================')

if __name__ == '__main__':
    try:
        main()
    except Exception:
        # no matter what error happen....just kill process 1 (which means, kill the vault server.)
        os.kill(1, signal.SIGTERM)