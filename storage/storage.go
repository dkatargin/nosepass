package storage

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"nosepass/encryption"
)

func StorePassword(dstPath string) error {
	publicKey, err := encryption.RSAGetPublicKey()
	if err != nil {
		return err
	}

	fmt.Print("Input password: ")
	binpass, err := terminal.ReadPassword(0)
	plainText := string(binpass)

	cipherText, err := encryption.RSAEncryptData(plainText, publicKey)
	if err != nil {
		return err
	}

	err = storeKey(dstPath, cipherText)
	if err != nil {
		return err
	}
	return nil
}

func GetPassword(dstPath string) (string, error) {
	key, err := getKey(dstPath)
	if err != nil {
		return "", err
	}

	privateKey, err := encryption.RSAGetPrivateKey()
	if err != nil {
		return "", err
	}
	plainText, err := encryption.RSADecryptData(key, privateKey)
	if err != nil {
		return "", err
	}

	return plainText, nil
}

func ListPassword() ([]string, error) {
	keys, err := listKeys()
	if err != nil {
		return nil, err
	}
	var pathList []string
	for _, secret := range keys {
		pathList = append(pathList, secret.path)
	}
	return pathList, nil
}

func DeletePassword(dstPath string) error {
	err := deleteKey(dstPath)
	return err
}
