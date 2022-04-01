package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	goutils "github.com/simonski/goutils"
	crypto "github.com/simonski/goutils/crypto"
)

// KPDB helper struct holds the data and keys
type KPDB struct {
	data               DB
	Filename           string
	PrivateKeyFilename string
}

// DB is the thing that we serialise to JSON
type DB struct {
	Entries map[string]DBEntry
}

// DBEntry represents the a single item in the DB
type DBEntry struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	LastUpdated time.Time `json:"lastUpdated"`
	Created     time.Time `json:"created"`
}

// NewKPDB constructor
func NewKPDB(filename string, privKey string) *KPDB {
	cdb := KPDB{}
	cdb.Load(filename, privKey)
	return &cdb
}

// Load populates the db with the file
func (cdb *KPDB) Load(filename string, privKey string) bool {
	cdb.Filename = goutils.EvaluateFilename(filename)
	cdb.PrivateKeyFilename = goutils.EvaluateFilename(privKey)

	if !goutils.FileExists(cdb.Filename) {
		db := DB{}
		db.Entries = make(map[string]DBEntry)
		cdb.data = db
	} else {
		jsonFile, err := os.Open(cdb.Filename)
		if err != nil {
			fmt.Printf("ERR %v\n", err)
			db := DB{}
			db.Entries = make(map[string]DBEntry)
			cdb.data = db
		} else {
			db := DB{}
			bytes, _ := ioutil.ReadAll(jsonFile)
			var data map[string]DBEntry
			json.Unmarshal(bytes, &data)
			db.Entries = data
			cdb.data = db
		}
	}

	return true
}

// Clear empties the db (without saving it)
func (cdb *KPDB) Clear() {
	cdb.data.Entries = make(map[string]DBEntry)
}

// Save writes the DB to disk
func (cdb *KPDB) Save() bool {
	data := cdb.data.Entries
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(cdb.Filename, file, 0644)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return true
}

// GetData returns the data map of all key
func (cdb *KPDB) GetData() DB {
	return cdb.data
}

// Get returns the (DBEntry, bool) indicating it exists (or not)
func (cdb *KPDB) Get(key string) (DBEntry, bool) {
	entry, exists := cdb.data.Entries[key]
	if exists {
		decValue := entry.Value
		decValue = cdb.Decrypt(entry.Value)
		entry.Value = decValue
	}
	return entry, exists
}

// Get returns the (DBEntry, bool) indicating it exists (or not)
func (cdb *KPDB) UpdateDescription(key string, description string) (DBEntry, bool) {
	entry, exists := cdb.Get(key)
	cdb.Put(entry.Key, entry.Value, description)
	if exists {
		entry.Description = description
	}
	return entry, exists
}

// Put stores (or replaces) the key/value pair
func (cdb *KPDB) Put(key string, value string, description string) {
	entry, exists := cdb.data.Entries[key]
	encValue := value
	encValue = cdb.Encrypt(value)
	if exists {
		if value != "" {
			entry.Value = encValue
		}
		entry.LastUpdated = time.Now()
		if description != "" {
			entry.Description = description
		}
		cdb.data.Entries[key] = entry
	} else {
		entry = DBEntry{Key: key, Value: encValue, Created: time.Now(), LastUpdated: time.Now()}
		if description != "" {
			entry.Description = description
		}
		cdb.data.Entries[key] = entry
	}
}

// Delete removes the key/value pair from the DB
func (cdb *KPDB) Delete(key string) {
	delete(cdb.data.Entries, key)
}

// Encrypt helper function encrypts with public key
func (cdb *KPDB) Encrypt(value string) string {
	encrypted, _ := crypto.EncryptWithPrivateKeyFilename(value, cdb.PrivateKeyFilename)
	return encrypted
}

// Decrypt helper function decrypts with private key
func (cdb *KPDB) Decrypt(value string) string {
	decrypted, _ := crypto.DecryptWithPrivateKeyFilename(value, cdb.PrivateKeyFilename)
	return decrypted
}
