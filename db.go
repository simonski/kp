package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type CrypticDB struct {
	data               DB
	Filename           string
	PublicKeyFilename  string
	PrivateKeyFilename string
}

type DB struct {
	Entries map[string]DBEntry
}

type DBEntry struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func NewCrypticDB(filename string, pubKey string, privKey string) *CrypticDB {
	db := CrypticDB{}
	db.Load(filename, pubKey, privKey)
	return &db
}

func evaluateFilename(filename string) string {
	home := os.Getenv("HOME")
	newname := strings.ReplaceAll(filename, "~", home)
	return newname
}

func (cdb *CrypticDB) Load(filename string, pubKey string, privKey string) bool {
	cdb.Filename = evaluateFilename(filename)
	cdb.PublicKeyFilename = evaluateFilename(pubKey)
	cdb.PrivateKeyFilename = evaluateFilename(privKey)
	jsonFile, err := os.Open(cdb.Filename)
	if err != nil {
		fmt.Printf("ERR %v\n", err)
		db := DB{}
		db.Entries = make(map[string]DBEntry)
		cdb.data = db
		// panic(err)
	} else {
		db := DB{}
		bytes, _ := ioutil.ReadAll(jsonFile)
		var data map[string]DBEntry
		json.Unmarshal(bytes, &data)
		db.Entries = data
		cdb.data = db
	}

	return true
}

func (cdb *CrypticDB) Clear() {
	cdb.data.Entries = make(map[string]DBEntry)
}

func (cdb *CrypticDB) Save() bool {
	data := cdb.data.Entries
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(cdb.Filename, file, 0644)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return true
}

func (cdb *CrypticDB) GetData() DB {
	return cdb.data
}

func (cdb *CrypticDB) Get(key string) (DBEntry, bool) {
	entry, exists := cdb.data.Entries[key]
	return entry, exists
}

func (cdb *CrypticDB) Put(key string, value string) {
	entry, exists := cdb.data.Entries[key]
	if exists {
		entry.Value = value
		cdb.data.Entries[key] = entry
	} else {
		entry = DBEntry{Key: key, Value: value}
		cdb.data.Entries[key] = entry
	}
}

func (db *CrypticDB) Delete(key string) {
	delete(db.data.Entries, key)
}

func (db *CrypticDB) Encrypt(value string) string {
	return value
}

func (db *CrypticDB) Decrypt(value string) string {
	return value
}
