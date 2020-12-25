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
	"nosepass/common"
	"os"
)

func RSAGenerateKeyPair() (*pem.Block, *pem.Block, error) {
	// Get passphrase
	fmt.Print("Input private key passphrase: \n")
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
	// Encrypt private key by passphrase
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
	configuration, err := common.Config()
	if err != nil {
		return nil, errors.New("wrong config")
	}
	// Check keydir
	if _, err = os.Stat(configuration.KeyDir); os.IsNotExist(err) {
		os.MkdirAll(configuration.KeyDir, 0700)
	}
	// Check keys
	privateKeyExist := false
	publicKeyExist := false
	if _, err := os.Stat(configuration.KeyDir + "private.pem"); !os.IsNotExist(err) {
		privateKeyExist = true
	}
	if _, err := os.Stat(configuration.KeyDir + "public.pem"); !os.IsNotExist(err) {
		publicKeyExist = true
	}
	// Generate keypair if not exist
	if !privateKeyExist {
		privblock, pubblock, err := RSAGenerateKeyPair()
		if err != nil {
			return nil, err
		}
		privatePem, err := os.Create(configuration.KeyDir + "private.pem")
		if err != nil {
			return nil, errors.New("error creation private.pem file")
		}
		err = pem.Encode(privatePem, privblock)
		if err != nil {
			return nil, errors.New("error encoding private.pem block")
		}
		publicPem, err := os.Create(configuration.KeyDir + "public.pem")
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

	// Read public key from file
	publicDat, err := ioutil.ReadFile(configuration.KeyDir + "public.pem")
	block, _ := pem.Decode(publicDat)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublickey, _ := pub.(*rsa.PublicKey)
	return rsaPublickey, nil
}

func RSAGetPrivateKey() (*rsa.PrivateKey, error) {
	configuration, err := common.Config()
	if err != nil {
		return nil, errors.New("wrong config")
	}
	// Check private key
	if _, err := os.Stat(configuration.KeyDir + "private.pem"); os.IsNotExist(err) {
		return nil, errors.New("private key not exist")
	}
	// Read private key from file
	privateDat, err := ioutil.ReadFile(configuration.KeyDir + "private.pem")
	if err != nil {
		return nil, err
	}
	fmt.Print("Input private key passphrase: \n")
	// Ask user passphrase
	binpass, err := terminal.ReadPassword(0)
	passphrase := string(binpass)
	// Decode pem block
	pemBlock, _ := pem.Decode(privateDat)
	if err != nil {
		return nil, err
	}
	// Decrypt pem block
	der, err := x509.DecryptPEMBlock(pemBlock, []byte(passphrase))
	if err != nil {
		return nil, err
	}
	// Parse private key
	privateKey, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return nil, err
	}
	return privateKey, err

}

func RSAEncryptData(value string, publicKey *rsa.PublicKey) (string, error) {
	// Encrypt data by public key
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(value), nil)
	if err != nil {
		return "", err
	}
	hexData := hex.EncodeToString(encryptedBytes)
	return hexData, err
}

func RSADecryptData(hexData string, privateKey *rsa.PrivateKey) (string, error) {
	// Decrypt data by private key
	ciphertext, err := hex.DecodeString(hexData)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	return string(decryptedBytes), err
}
