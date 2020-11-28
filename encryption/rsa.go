package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"nosepass/storage"
	"os"
)

func RSAGenerateKeyPair() (*pem.Block, *pem.Block, error) {
	// Get passphrase
	fmt.Print("Input password: ")
	binpass, err := terminal.ReadPassword(0)
	passphrase := string(binpass)
	if err != nil {
		fmt.Println("Error write password: " + string(passphrase))
	}
	if passphrase == "" {
		return nil, nil, errors.New("empty passphrase")
	}

	// Private Key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	}

	privateKeyBlock, err = x509.EncryptPEMBlock(rand.Reader, privateKeyBlock.Type, privateKeyBlock.Bytes, []byte(passphrase), x509.PEMCipherAES256)

	// Public Key
	publickey := &privatekey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return nil, nil, errors.New("error public key marshal")
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	return privateKeyBlock, publicKeyBlock, nil
}

func RSAGetPublicKey() (*pem.Block, error) {
	var publickey *pem.Block
	configuration, err := storage.Config()
	if err != nil {
		return nil, errors.New("wrong config")
	}
	// check keydir
	if _, err := os.Stat(configuration.keydir); os.IsNotExist(err) {
		return nil, errors.New("keydir not exist")
	}
	// check keys
	privateKeyExist := false
	publicKeyExist := false
	if _, err := os.Stat(configuration.keydir + "private.pem"); !os.IsNotExist(err) {
		privateKeyExist = true
	}
	if _, err := os.Stat(configuration.keydir + "public.pem"); !os.IsNotExist(err) {
		publicKeyExist = true
	}
	// Generate keypair if not exist
	if !privateKeyExist {
		privblock, pubblock, err := RSAGenerateKeyPair()
		if err != nil {
			return nil, err
		}
		privatePem, err := os.Create(configuration.keydir + "private.pem")
		if err != nil {
			return nil, errors.New("error creation private.pem file")
		}
		err = pem.Encode(privatePem, privblock)
		if err != nil {
			return nil, errors.New("error encoding private.pem block")
		}
		publicPem, err := os.Create(configuration.keydir + "public.pem")
		if err != nil {
			return nil, errors.New("error creation public.pem file")
		}
		err = pem.Encode(publicPem, pubblock)
		if err != nil {
			return nil, errors.New("error encoding public.pem block")
		}
		publicKeyExist = true
		publickey = pubblock
	}
	// Check public key
	if !publicKeyExist {
		return nil, errors.New("public key not exist")
	}

	publicDat, err := ioutil.ReadFile(configuration.keydir + "public.pem")
	block, _ := pem.Decode(publicDat)
	return block, nil

}

func RSAEncryptData(value string, publicKey rsa.PublicKey) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, []byte("super secret message"), nil)
	if err != nil {
		return nil, err
	}
}
