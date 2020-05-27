package internal

import (
	"github.com/99designs/keyring"
	"github.com/KMConner/kyodai-go/kulasis"
)

const (
	appName    = "KYODAI_GO"
	accountKey = "Account"
	tokenKey   = "Token"
)

func openStore() (keyring.Keyring, error) {
	ring, err := keyring.Open(keyring.Config{
		AllowedBackends: []keyring.BackendType{keyring.WinCredBackend, keyring.KeychainBackend},
		ServiceName:     appName,
	})

	return ring, err
}

func Store(info kulasis.Info) error {
	ring, err := openStore()
	if err != nil {
		return err
	}
	err = ring.Set(keyring.Item{
		Key:  accountKey,
		Data: []byte(info.Account),
	})
	if err != nil {
		return err
	}

	err = ring.Set(keyring.Item{
		Key:  tokenKey,
		Data: []byte(info.AccessToken),
	})
	return err
}

func Load() (*kulasis.Info, error) {
	ring, err := openStore()
	if err != nil {
		return nil, err
	}
	info := kulasis.Info{}
	account, err := ring.Get(accountKey)
	if err != nil {
		return nil, err
	}
	info.Account = string(account.Data)

	token, err := ring.Get(tokenKey)
	if err != nil {
		return nil, err
	}
	info.AccessToken = string(token.Data)

	return &info, nil
}
