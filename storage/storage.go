package storage

import (
	"nosepass/encryption"
)

func StorePassword(dstPath string) error {
	pem, err := encryption.RSAGetPublicKey()
	if err != nil {
		return err
	}

	return nil
}
