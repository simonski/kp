package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
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
	Version string             `json:"version"`
	Entries map[string]DBEntry `json:"entries"`
}

// DBEntry represents the a single item in the DB
type DBEntry struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	Notes       string    `json:"notes"`
	Username    string    `json:"username"`
	Url         string    `json:"url"`
	Type        string    `json:"type"`
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
			bytes, _ := io.ReadAll(jsonFile)

			// let's see what sort of DB this is
			// pre-version was a map of Entries
			// post-version should be a DB
			// we can test by trying to load the DB directly.
			json.Unmarshal(bytes, &db)
			if db.Version == "" {
				var data map[string]DBEntry
				json.Unmarshal(bytes, &data)

				db.Entries = data
				db.Version = DB_VERSION
				cdb.data = db
				cdb.Save()
			} else {
				// then it has a schema version
				// really we should now check and upgrade, e.g.
				// if db.Version != DB_VERSION
				// upgrade()
				cdb.data = db
			}

			for k, v := range cdb.data.Entries {
				if strings.TrimSpace(v.Key) == "" {
					delete(cdb.data.Entries, v.Key)
					if strings.TrimSpace(k) == "" {
						k = uuid.New().String()
					}
					v.Key = k
					cdb.data.Entries[k] = v
				}
			}

		}
	}

	return true
}

func (cdb *KPDB) GetEntriesSortedByUpdatedThenKey() []DBEntry {

	entries := make([]DBEntry, 0)
	for _, e := range cdb.data.Entries {
		entries = append(entries, e)
	}

	sort.SliceStable(entries, func(a int, b int) bool {
		entryA := entries[a]
		entryB := entries[b]
		keyA := strings.ToLower(entryA.Key)
		keyB := strings.ToLower(entryB.Key)
		return keyA < keyB
	})

	return entries

}

// Clear empties the db (without saving it)
func (cdb *KPDB) Clear() {
	cdb.data.Entries = make(map[string]DBEntry)
}

// Save writes the DB to disk
func (cdb *KPDB) Save() bool {
	data := cdb.data
	file, _ := json.MarshalIndent(data, "", " ")
	err := os.WriteFile(cdb.Filename, file, 0644)
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
func (cdb *KPDB) GetDecrypted(key string) (DBEntry, bool) {
	entry, exists := cdb.data.Entries[key]
	if exists {
		decValue, _ := cdb.Decrypt(entry.Value)
		entry.Value = decValue
	}
	return entry, exists
}

// Get returns the (DBEntry, bool) indicating it exists (or not)
func (cdb *KPDB) UpdateDescription(key string, description string) (DBEntry, bool) {
	entry, exists := cdb.GetDecrypted(key)
	entry.Description = description
	cdb.Put(entry)
	return entry, exists
}

// Put stores (or replaces) the key/value pair
func (cdb *KPDB) Put(entry_in DBEntry) {
	entry, exists := cdb.data.Entries[entry_in.Key]
	encValue, _ := cdb.Encrypt(entry_in.Value)
	entry_in.Value = encValue
	if !exists {
		entry_in.Created = entry.Created
	} else {
		entry_in.Created = time.Now()
	}
	entry_in.LastUpdated = time.Now()
	cdb.data.Entries[entry_in.Key] = entry_in
}

// Delete removes the key/value pair from the DB
func (cdb *KPDB) Delete(key string) {
	delete(cdb.data.Entries, key)
}

// Encrypt helper function encrypts with public key
func (cdb *KPDB) Encrypt(value string) (string, error) {
	return crypto.EncryptWithPrivateKeyFilename(value, cdb.PrivateKeyFilename)
}

// Decrypt helper function decrypts with private key
func (cdb *KPDB) Decrypt(value string) (string, error) {
	return crypto.DecryptWithPrivateKeyFilename(value, cdb.PrivateKeyFilename)
}
