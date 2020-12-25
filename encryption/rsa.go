package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
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

func RSAGetPublicKey() (*rsa.PublicKey, error) {
	//var publickey *pem.Block
	configuration, err := Config()
	if err != nil {
		return nil, errors.New("wrong config")
	}
	// check keydir
	if _, err := os.Stat(configuration.Keydir); os.IsNotExist(err) {
		return nil, errors.New("keydir not exist")
	}
	// check keys
	privateKeyExist := false
	publicKeyExist := false
	if _, err := os.Stat(configuration.Keydir + "private.pem"); !os.IsNotExist(err) {
		privateKeyExist = true
	}
	if _, err := os.Stat(configuration.Keydir + "public.pem"); !os.IsNotExist(err) {
		publicKeyExist = true
	}
	// Generate keypair if not exist
	if !privateKeyExist {
		privblock, pubblock, err := RSAGenerateKeyPair()
		if err != nil {
			return nil, err
		}
		privatePem, err := os.Create(configuration.Keydir + "private.pem")
		if err != nil {
			return nil, errors.New("error creation private.pem file")
		}
		err = pem.Encode(privatePem, privblock)
		if err != nil {
			return nil, errors.New("error encoding private.pem block")
		}
		publicPem, err := os.Create(configuration.Keydir + "public.pem")
		if err != nil {
			return nil, errors.New("error creation public.pem file")
		}
		err = pem.Encode(publicPem, pubblock)
		if err != nil {
			return nil, errors.New("error encoding public.pem block")
		}
		publicKeyExist = true
		//publickey = pubblock
	}
	// Check public key
	if !publicKeyExist {
		return nil, errors.New("public key not exist")
	}

	publicDat, err := ioutil.ReadFile(configuration.Keydir + "public.pem")
	block, _ := pem.Decode(publicDat)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublickey, _ := pub.(*rsa.PublicKey)
	return rsaPublickey, nil
}

func RSAGetPrivateKey() (*rsa.PrivateKey, error) {
	configuration, err := Config()
	if err != nil {
		return nil, errors.New("wrong config")
	}
	if _, err := os.Stat(configuration.Keydir + "private.pem"); os.IsNotExist(err) {
		return nil, errors.New("private key not exist")
	}
	privateDat, err := ioutil.ReadFile(configuration.Keydir + "private.pem")
	if err != nil {
		return nil, err
	}
	fmt.Print("Input private key password: ")
	binpass, err := terminal.ReadPassword(0)
	passphrase := string(binpass)
	pemBlock, _ := pem.Decode(privateDat)
	if err != nil {
		return nil, err
	}
	der, err := x509.DecryptPEMBlock(pemBlock, []byte(passphrase))
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return nil, err
	}
	return privateKey, err

}

func RSAEncryptData(value string, publicKey *rsa.PublicKey) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(value), nil)
	if err != nil {
		return "", err
	}
	hexData := hex.EncodeToString(encryptedBytes)
	return hexData, err
}

func RSADecryptData(hexData string, privateKey *rsa.PrivateKey) (string, error) {
	ciphertext, err := hex.DecodeString(hexData)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	return string(decryptedBytes), err
}
