package main

import (
	"context"
	"fmt"
	"encoding/hex"

	schnorrkel "github.com/ChainSafe/go-schnorrkel"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type AccountJSON struct {
	PublicKey string `json:"publicKey"`
	PrivateKey [32]byte `json:"privateKey"`
}

func accountPaths(b *PluginBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: QualifiedPath("accounts/" + framework.GenericNameRegex("name")),
			HelpSynopsis: "Generates an account",
			HelpDescription: "Generates an account",
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathAccountsCreate,
				logical.ReadOperation: b.pathAccountsRead,
			},
		},
		{
			Pattern: QualifiedPath("accounts/" + framework.GenericNameRegex("name") + "/sign"),
			HelpSynopsis: "Sign a message",
			HelpDescription: "Sign a message",
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"message": {
					Type:        framework.TypeString,
					Description: "Message to sign.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSignMessage,
				logical.UpdateOperation: b.pathSignMessage,
			},
		},
	}
}

func (b *PluginBackend) updateAccount(ctx context.Context, req *logical.Request, name string, accountJSON *AccountJSON) error {
	path := QualifiedPath(fmt.Sprintf("accounts/%s", name))

	entry, err := logical.StorageEntryJSON(path, accountJSON)
	if err != nil {
		return err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return err
	}
	return nil
}

func readAccount(ctx context.Context, req *logical.Request, name string) (*AccountJSON, error) {
	path := QualifiedPath(fmt.Sprintf("accounts/%s", name))

	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var accountJSON AccountJSON
	err = entry.DecodeJSON(&accountJSON)

	if entry == nil {
		return nil, fmt.Errorf("failed to deserialize account at %s", path)
	}
	return &accountJSON, nil
}

func (b *PluginBackend) pathAccountsCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	priv, pub, err := schnorrkel.GenerateKeypair()
	if err != nil {
		return nil, err
	}

	pubEncoded := pub.Encode()
	pubHex := "0x"+hex.EncodeToString(pubEncoded[:])

	accountJSON := &AccountJSON{
		PublicKey: pubHex,
		PrivateKey: priv.Encode(),
	}

	err = b.updateAccount(ctx, req, name, accountJSON)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"publicKey": pubHex,
		},
	}, nil
}

func (b *PluginBackend) pathAccountsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	accountJSON, err := readAccount(ctx, req, name)

	if err != nil || accountJSON == nil {
		return nil, fmt.Errorf("Error reading account")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"publicKey": accountJSON.PublicKey,
		},
	}, nil
}

func (b *PluginBackend) pathSignMessage(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	message := data.Get("message").(string)
	name := data.Get("name").(string)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}

	msg := []byte(message)
	var p [32]byte
	copy(p[:], accountJSON.PrivateKey[:])
	secretKey := &schnorrkel.SecretKey{}
	secretKey.Decode(p)
	signingCtx := []byte("substrate")
	signingTranscript := schnorrkel.NewSigningContext(signingCtx, msg)
	sig, err := secretKey.Sign(signingTranscript)
	if err != nil {
		return nil, err
	}

	sigEncoded := sig.Encode()
	sigHex := "0x"+hex.EncodeToString(sigEncoded[:])

	return &logical.Response{
		Data: map[string]interface{}{
			"message": message,
			"publicKey": accountJSON.PublicKey,
			"signature": sigHex,
		},
	}, nil
}

func pathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}
