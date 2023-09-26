package main

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func FactoryType(backendType logical.BackendType) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b, err := Backend(conf)
		if err != nil {
			return nil, err
		}
		b.BackendType = backendType
		if err = b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

func Backend(conf *logical.BackendConfig) (*PluginBackend, error) {
	var b PluginBackend
	b.Backend = &framework.Backend{
		Help: "",
		Paths: framework.PathAppend(
			accountPaths(&b),
		),
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}
	return &b, nil
}

type PluginBackend struct {
	*framework.Backend
}

func QualifiedPath(subpath string) string {
	return subpath
}
