# Verification and Vault Plugin Repository

This repository contains two main components: a Verification Tool for signing and verifying messages, and a Vault Plugin for account management and message signing within HashiCorp Vault.

## Getting Started

Make sure to have [Node.js](https://nodejs.org/) and [Go](https://golang.org/) installed on your machine before proceeding.

### 1. Verification Tool

This tool allows you to sign a test message using either ECDSA or SR25519 cryptographic algorithms, and then verify the signed message.

#### Installation

```bash
npm install
```

#### Usage

First, navigate to either the `ecdsa` or `sr25519` directory, then run the following command to sign a test message:

```bash
go run sign.go "test message"
```

Take note of the output, then run the following command to verify the signature:

```bash
node verify.js [output-from-previous-step]
```

In case the above command fails, use:

```bash
node -r esm verify.js [output-from-previous-step]
```

### 2. Vault Plugin

This plugin extends HashiCorp Vault to provide account management and message signing capabilities.

#### Installation

First, install Vault using the following commands:

```bash
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install vault
```

Then, navigate to the `vault` directory and run the setup script to build the plugin and start a local Vault server in development mode:

```bash
cd vault
sh setup.sh
```

#### Usage

Set the address for the local Vault server:

```bash
export VAULT_ADDR='http://127.0.0.1:8200'
```

To create a new account (which generates an internal key pair and returns the public key):

```bash
vault write avn-vault/accounts/test_account name="test_account"
```

To sign a message using the newly created account:

```bash
vault write avn-vault/accounts/test_account/sign message="test message"
```

To retrieve the public key of an account:

```bash
vault read avn-vault/accounts/test_account
```

Utilize the Verification Tool to verify the signed message using the returned public key.
