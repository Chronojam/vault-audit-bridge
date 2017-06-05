#!/bin/bash

docker run -d --name=vault --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' --net=host vault
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=myroot

## Wait for vault to come up, lazy
sleep 2

vault audit-enable socket address="127.0.0.1:3333" socket_type="tcp"
