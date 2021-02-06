package main

import (
	"fmt"
	"os"
	"sort"

	clipboard "github.com/atotto/clipboard"
	goutils "github.com/simonski/goutils"
)

func main() {
	cli := goutils.NewCLI(os.Args)
	command := cli.GetCommand()
	if command == "help" {
		fmt.Printf("ssh-keygen -m pem -f ~/.ssh/id_rsa\n")
		fmt.Printf("ssh-keygen -f ~/.ssh/id_rsa.pub -e -m pem > ~/.ssh/id_rsa.pem\n")
	} else if isInfo(command) {
		DoInfo(cli)
	} else if isVersion(command) {
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

func isInfo(command string) bool {
	return command == "info"
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

func LoadDB() *CrypticDB {
	filename := goutils.GetEnvOrDefault(CRYPTIC_FILE, "~/.Crypticfile")
	pubKey := goutils.GetEnvOrDefault(CRYPTIC_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privKey := goutils.GetEnvOrDefault(CRYPTIC_PRIVATE_KEY, "~/.ssh/id_rsa")
	encryptionEnabled := goutils.GetEnvOrDefault(CRYPTIC_ENCRYPTION_ENABLED, "1") == "1"
	db := NewCrypticDB(filename, pubKey, privKey, encryptionEnabled)
	return db
}

func DoInfo(cli *goutils.CLI) {
	filename := goutils.GetEnvOrDefault(CRYPTIC_FILE, "~/.Crypticfile")
	publicKey := goutils.GetEnvOrDefault(CRYPTIC_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privateKey := goutils.GetEnvOrDefault(CRYPTIC_PRIVATE_KEY, "~/.ssh/id_rsa")

	filenameExists := goutils.FileExists(goutils.EvaluateFilename(filename))
	publicKeyExists := goutils.FileExists(goutils.EvaluateFilename(publicKey))
	privateKeyExists := goutils.FileExists(goutils.EvaluateFilename(privateKey))

	fmt.Printf("%v          =%v, exists=%v\n", CRYPTIC_FILE, filename, filenameExists)
	fmt.Printf("%v =%v, exists=%v\n", CRYPTIC_PUBLIC_KEY, publicKey, publicKeyExists)
	fmt.Printf("%v =%v, exists=%v\n", CRYPTIC_PRIVATE_KEY, privateKey, privateKeyExists)
	// fmt.Printf("%v =%v, exists=%v\n", CRYPTIC_ENCRYPTION_ENABLED, privateKeyExists)
}

func DoGet(cli *goutils.CLI) {
	key := cli.GetCommand()
	db := LoadDB()
	entry, exists := db.Get(key)
	if exists {
		value := entry.Value
		clipboard.WriteAll(value)
		if cli.IndexOf("-stdout") > -1 {
			fmt.Printf("%v\n", value)
		}
	} else {
		fmt.Printf("'%v' does not exist.\n", key)
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
