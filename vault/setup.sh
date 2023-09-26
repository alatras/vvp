#!/bin/bash
go build -o plugin/avn-vault
sleep 5
gnome-terminal -- vault server -dev -dev-root-token-id=root -dev-plugin-dir=./plugin
sleep 3
clear
export VAULT_ADDR='http://127.0.0.1:8200'
vault login root
vault secrets enable -path=avn-vault avn-vault