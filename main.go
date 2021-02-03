package main

import (
	"fmt"
	"os"
	"sort"

	clipboard "github.com/atotto/clipboard"
	goutils "github.com/simonski/goutils"
)

const VERSION = "1.0.0"
const GLOBAL_USAGE = `cryptic is a tool for using key/pairs.

Usage:

    cryptic <key | command> <value>

The commands are:

    ls                  list keys
    rm [key]            remove key "key"
    key                 get the value of "key"
    key value           overwrite the value of "key"

    clear               remove all values
    version             print application version"

`

func main() {
	cli := goutils.NewCLI(os.Args)
	command := cli.GetCommand()
	if command == "test" {
		db := LoadDB()
		value := cli.GetStringOrDie(command)
		value_enc := db.Encrypt(value)
		value_dec := db.Decrypt(value_enc)
		fmt.Printf("Encrypt '%v' = '%v', decrypt = '%v'\n", value, value_enc, value_dec)
		os.Exit(0)
	}
	if isVersion(command) {
		DoVersion(cli)
	} else if isClear(command) {
		DoClear(cli)
	} else if isList(command) {
		DoList(cli)
	} else if isGet(command, cli) {
		DoGet(cli)
	} else if isPut(command, cli) {
		DoPut(cli)
	} else if isDelete(command) {
		DoDelete(cli)
	} else {
		DoUsage(cli)
	}
}

func isDelete(command string) bool {
	return command == "rm"
}

func isVersion(command string) bool {
	return command == "version"
}

func isList(command string) bool {
	return command == "ls"
}

func isClear(command string) bool {
	return command == "clear"
}

// A 'get' is basically not a a list, delete or a put
func isGet(command string, cli *goutils.CLI) bool {
	return command != "" && !isPut(command, cli) && !isDelete(command) && !isList(command)
}

func isPut(command string, cli *goutils.CLI) bool {
	if isDelete(command) || isList(command) {
		return false
	}
	value := cli.GetStringOrDefault(command, "")
	if value == "" {
		// can't be a put as there is no value
		return false
	}
	return true

}

func GetEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	} else {
		return defaultValue
	}
}

func LoadDB() *CrypticDB {
	filename := GetEnvOrDefault("CRYPTIC_DB", "~/.crypticdb.json")
	pubKey := GetEnvOrDefault("CRYPTIC_PUB_KEY", "~/.ssh/id_rsa.pub")
	privKey := GetEnvOrDefault("CRYPTIC_PRIV_KEY", "~/.ssh/id_rsa")
	db := NewCrypticDB(filename, pubKey, privKey)
	return db
}

func DoGet(cli *goutils.CLI) {
	key := cli.GetCommand()
	db := LoadDB()
	entry, exists := db.Get(key)
	if exists {
		value := entry.Value
		clipboard.WriteAll(value)
	} else {
		fmt.Printf("'%v' does not exist.", key)
		os.Exit(1)
	}
}

func DoPut(cli *goutils.CLI) {
	db := LoadDB()
	key := cli.GetCommand()
	value := cli.GetStringOrDefault(key, "")
	db.Put(key, value)
	db.Save()
}

func DoClear(cli *goutils.CLI) {
	db := LoadDB()
	db.Clear()
	db.Save()
}

func DoList(cli *goutils.CLI) {
	db := LoadDB()
	data := db.GetData()
	if len(data.Entries) == 0 {
		fmt.Printf("DB is empty.\n")
	} else {
		maxLength := 0
		for key, _ := range data.Entries {
			maxLength = goutils.Max(len(key), maxLength)
		}

		keys := make([]string, 0)
		for key, _ := range data.Entries {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf("%v\n", key)
		}
	}
}

func DoDelete(cli *goutils.CLI) {
	command := cli.GetCommand()
	key := cli.GetStringOrDefault(command, "")
	if key == "" {
		USAGE := "cryptic rm [key]"
		fmt.Printf("%v\n", USAGE)
	}
	db := LoadDB()
	db.Delete(key)
	db.Save()

}

func DoUsage(cli *goutils.CLI) {
	fmt.Printf(GLOBAL_USAGE)
}

func DoVersion(cli *goutils.CLI) {
	fmt.Printf("%v\n", VERSION)
}
