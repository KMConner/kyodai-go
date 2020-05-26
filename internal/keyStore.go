package internal

import (
	"github.com/99designs/keyring"
	"github.com/KMConner/kyodai-go/kulasis"
)

const appName = "KYODAI_GO"

func openStore() (keyring.Keyring, error) {
	ring, err := keyring.Open(keyring.Config{
		AllowedBackends: []keyring.BackendType{keyring.WinCredBackend},
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
		Key:  "Account",
		Data: []byte(info.Account),
	})
	if err != nil {
		return err
	}

	err = ring.Set(keyring.Item{
		Key:  "Token",
		Data: []byte(info.AccessToken),
	})
	return err
}
