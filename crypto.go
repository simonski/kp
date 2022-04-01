package main

/*
crypto/rsa
The func GenerateKey(random io.Reader, bits int) (*PrivateKey, error) function generates an RSA keypair of the given bit size using the random source random (for example, crypto/rand.Reader, discussed below).

crypto/rand
The var Reader io.Reader struct is a global, shared instance of a cryptographically secure random number generator. On Linux, Reader uses getrandom(2) if available, /dev/urandom otherwise.

crypto/x509
The func MarshalPKCS1PrivateKey(key *rsa.PrivateKey) []byte function converts an RSA private key to PKCS #1, ASN.1 DER form. This kind of key is commonly encoded in PEM blocks of type “RSA PRIVATE KEY”.

The func MarshalPKCS1PublicKey(key *rsa.PublicKey) []byte function converts an RSA public key to PKCS #1, ASN.1 DER form. This kind of key is commonly encoded in PEM blocks of type “RSA PUBLIC KEY”.

encoding/pem
The func Encode(out io.Writer, b *Block) error function writes the PEM encoding of b to out.



*/
import (
	"crypto/rsa"
	b64 "encoding/base64"
	"errors"
	"io/ioutil"

	goutils "github.com/simonski/goutils"

	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func Encrypt(value string, privateKeyFilename string) (string, error) {
	return EncryptWithPrivateKeyFilename(value, privateKeyFilename)
}

// Decrypt helper function decrypts with private key
func Decrypt(value string, privateKeyFilename string) (string, error) {
	f := goutils.EvaluateFilename(privateKeyFilename)
	uDec, _ := b64.StdEncoding.DecodeString(value)
	privateKey, err := LoadPrivateKey(f)
	if err != nil {
		return "", err
	}
	bytes := []byte(uDec)
	decrypted := DecryptWithPrivateKey(bytes, privateKey)
	s := string(decrypted)
	return s, nil
}

// LoadPrivateKey loads the filename to a *rsa.PrivateKey
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	f := goutils.EvaluateFilename(filename)
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return BytesToPrivateKey(bytes)
}

// LoadPublicKey loads the filename to a *rsa.PublicKey
func LoadPublicKey(filename string) *rsa.PublicKey {
	f := goutils.EvaluateFilename(filename)
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return BytesToPublicKey(bytes)
}

// DoVerify performs verification of ~/.KPfile, encryption/decryption using
// specified keys
func Verify(cli *goutils.CLI, printFailuresToStdOut bool) bool {
	overallValid := true

	keyFilename := cli.GetFileExistsOrDie("-key")
	keyExists := goutils.FileExists(goutils.EvaluateFilename(keyFilename))

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("%v   %v\n", "key", keyFilename))

	if !keyExists {
		line := fmt.Sprintf("key file '%v' does not exist.\n", keyFilename)
		messages = append(messages, line)
		overallValid = false
	}

	if keyExists {
		// try to encrypt/decrypt something
		plain := "Hello, World"
		fmt.Printf("encrypting %v with key %v\n", plain, keyFilename)
		encrypted, err := EncryptWithPrivateKeyFilename(plain, keyFilename)
		if err != nil {
			line := fmt.Sprintf("Cannot encrypt: %v.\n", err)
			messages = append(messages, line)
		}

		decrypted, err := DecryptWithPrivateKeyFilename(encrypted, keyFilename)
		if err != nil {
			line := fmt.Sprintf("Cannot decrypt: %v.\n", err)
			messages = append(messages, line)
		}
		if plain != decrypted {
			line := "Encrypt/Decrypt not working.\n"
			messages = append(messages, line)
			overallValid = false
		}

	} else {

		messages = append(messages, "\nPublic/private keys do not exist, try the following\n\n")
		line := "    ssh-keygen -b 2048 -t rsa -m pem -f crypto/id_rsa\n"
		messages = append(messages, line)
		line = "    <or>\n"
		messages = append(messages, line)
		line = "    ssh-keygen -b 2048 -t rsa -m pkcs8 -f crypto/id_rsa\n"
		messages = append(messages, line)
		line = "    ssh-keygen -f crypto/id_rsa.pub -e -m pem > crypto/id_rsa.pem\n\n"
		messages = append(messages, line)
	}

	if printFailuresToStdOut {
		for _, line := range messages {
			fmt.Print(line)
		}
	}

	if overallValid {
		if printFailuresToStdOut {
			fmt.Printf("Verify : OK.\n")
		}
	}

	return overallValid
}

func EncryptWithPublicKeyFilename(value string, publicKeyFilename string) string {
	f := goutils.EvaluateFilename(publicKeyFilename)
	publicKey := LoadPublicKey(f)
	bytes := []byte(value)
	encrypted := EncryptWithPublicKey(bytes, publicKey)
	s := b64.StdEncoding.EncodeToString(encrypted)
	return s
}

func EncryptWithPrivateKeyFilename(value string, privateKeyFilename string) (string, error) {
	f := goutils.EvaluateFilename(privateKeyFilename)
	publicKey, err := LoadPublicFromPrivateKey(f)
	if err != nil {
		return "", err
	}
	bytes := []byte(value)
	encrypted := EncryptWithPublicKey(bytes, publicKey)
	s := b64.StdEncoding.EncodeToString(encrypted)
	return s, nil
}

// Decrypt helper function decrypts with private key
func DecryptWithPrivateKeyFilename(value string, privateKeyFilename string) (string, error) {
	f := goutils.EvaluateFilename(privateKeyFilename)
	uDec, _ := b64.StdEncoding.DecodeString(value)
	privateKey, err := LoadPrivateKey(f)
	if err != nil {
		return "", err
	}
	bytes := []byte(uDec)
	decrypted := DecryptWithPrivateKey(bytes, privateKey)
	s := string(decrypted)
	return s, nil
}

// LoadPublicKey loads the filename to a *rsa.PublicKey
func LoadPublicFromPrivateKey(filename string) (*rsa.PublicKey, error) {
	f := goutils.EvaluateFilename(filename)
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	privateKey, err := BytesToPrivateKey(bytes)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey
	return &publicKey, nil

}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	if block == nil {
		fmt.Printf("Error, block is nil.\n")
		os.Exit(1)
	}
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		// fmt.Println("decrypting the pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
		// } else {
		// 	fmt.Println("not decrypting the pem block")
	}

	var parsedKey interface{}

	if parsedKey, err = x509.ParsePKCS1PrivateKey(b); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(b); err != nil {
			return nil, errors.New("Cannot load private key - neither PKCS1 not PKCS8.")
		}
	}
	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Cannot load private key.")
	}

	return privateKey, nil
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	if block == nil {
		fmt.Printf("Error, block is nil.\n")
		os.Exit(1)
	}
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			panic(err)
		}
	}
	ifc, err := x509.ParsePKCS1PublicKey(b)

	if err != nil {
		panic(err)
	}
	return ifc
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		panic(err) //log.Error(err)
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		panic(err) //log.Error(err)
	}
	return plaintext
}

func Create() {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Cannot generate RSA key\n")
		os.Exit(1)
	}
	// publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create("private.pem")
	if err != nil {
		fmt.Printf("error when create private.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encode private pem: %s \n", err)
		os.Exit(1)
	}

	// // dump public key to file
	// publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	// if err != nil {
	// 	fmt.Printf("error when dumping publickey: %s \n", err)
	// 	os.Exit(1)
	// }
	// publicKeyBlock := &pem.Block{
	// 	Type:  "PUBLIC KEY",
	// 	Bytes: publicKeyBytes,
	// }
	// publicPem, err := os.Create("public.pem")
	// if err != nil {
	// 	fmt.Printf("error when create public.pem: %s \n", err)
	// 	os.Exit(1)
	// }
	// err = pem.Encode(publicPem, publicKeyBlock)
	// if err != nil {
	// 	fmt.Printf("error when encode public pem: %s \n", err)
	// 	os.Exit(1)
	// }
}
