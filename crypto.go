package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// LoadPrivateKey loads the filename to a *rsa.PrivateKey
func LoadPrivateKey(filename string) *rsa.PrivateKey {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return BytesToPrivateKey(bytes)
}

// LoadPublicKey loads the filename to a *rsa.PublicKey
func LoadPublicKey(filename string) *rsa.PublicKey {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return BytesToPublicKey(bytes)
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	// fmt.Printf("block.Type=%v\n", block.Type)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			panic(err) //log.Error(err)
		}

	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(b); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(b); err != nil {
			fmt.Printf("neither 1 nor 8\n")
		}
	}
	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		fmt.Printf("cannot\n")
	}

	// key, err := x509.ParsePKCS8PrivateKey(b)
	// fmt.Printf("key=%v\n", key)
	// fmt.Printf("err=%v\n", err)
	// if err != nil {
	// 	panic(err) //log.Error(err)
	// }
	return privateKey
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	// fmt.Printf("BytesToPublicKey, pub=%v\n", pub)
	block, _ := pem.Decode(pub)
	// fmt.Printf("BytesToPublicKey, pub=%v\n", pub)
	if block == nil {
		fmt.Printf("Error, block is nill.\n")
		os.Exit(1)
	}
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			panic(err) //log.Error(err)
		}
	}
	ifc, err := x509.ParsePKCS1PublicKey(b)

	// ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		panic(err) //log.Error(err)
	}
	// key, ok := ifc.(*rsa.PublicKey)
	// if !ok {
	// 	panic(err) //log.Error(err)
	// }
	// return key
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
