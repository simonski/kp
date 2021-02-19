package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	clipboard "github.com/atotto/clipboard"
	goutils "github.com/simonski/goutils"

	terminal "golang.org/x/crypto/ssh/terminal"
)

func main() {
	cli := goutils.NewCLI(os.Args)
	command := cli.GetCommand()
	if command == "help" {
		DoUsage(cli)
	} else if isInfo(command) {
		DoInfo(cli)
	} else if isVerify(command) {
		DoVerify(cli, false)
	} else if isVersion(command) {
		DoVersion(cli)
	} else if isClear(command) {
		DoClear(cli)
	} else if isList(command) {
		DoList(cli)
	} else if isPut(command, cli) {
		// DoVerify(cli, true)
		DoPut(cli)
	} else if isGet(command, cli) {
		// DoVerify(cli, true)
		DoGet(cli)
	} else if isDelete(command) {
		DoDelete(cli)
	} else if command != "" {
		fmt.Printf("cryptic %v: unknown command\n", command)
		fmt.Printf("Run 'cryptic help' for usage.\n")
		os.Exit(1)
	} else {
		DoUsage(cli)
	}
}

func isVerify(command string) bool {
	return command == "verify"
}

func isInfo(command string) bool {
	return command == "info"
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
	return command == "get"
}

func isPut(command string, cli *goutils.CLI) bool {
	return command == "put"
}

func LoadDB() *CrypticDB {
	filename := goutils.GetEnvOrDefault(CRYPTIC_FILE, "~/.Crypticfile")
	pubKey := goutils.GetEnvOrDefault(CRYPTIC_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privKey := goutils.GetEnvOrDefault(CRYPTIC_PRIVATE_KEY, "~/.ssh/id_rsa")
	encryptionEnabled := goutils.GetEnvOrDefault(CRYPTIC_ENCRYPTION_ENABLED, "1") == "1"
	db := NewCrypticDB(filename, pubKey, privKey, encryptionEnabled)
	return db
}

// DoVerify performs verification of ~/.Crypticfile, encryption/decryption using
// specified keys
func DoVerify(cli *goutils.CLI, failOnError bool) {
	overallValid := true
	filename := goutils.GetEnvOrDefault(CRYPTIC_FILE, "~/.Crypticfile")
	publicKey := goutils.GetEnvOrDefault(CRYPTIC_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privateKey := goutils.GetEnvOrDefault(CRYPTIC_PRIVATE_KEY, "~/.ssh/id_rsa")

	filenameExists := goutils.FileExists(goutils.EvaluateFilename(filename))
	publicKeyExists := goutils.FileExists(goutils.EvaluateFilename(publicKey))
	privateKeyExists := goutils.FileExists(goutils.EvaluateFilename(privateKey))

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("%v          %v\n", CRYPTIC_FILE, filename))
	messages = append(messages, fmt.Sprintf("%v    %v\n", CRYPTIC_PUBLIC_KEY, publicKey))
	messages = append(messages, fmt.Sprintf("%v   %v\n", CRYPTIC_PRIVATE_KEY, privateKey))

	overallValid = filenameExists && publicKeyExists && privateKeyExists
	if !filenameExists {
		line := fmt.Sprintf("Crypticfile '%v' does not exist.\n", filename)
		messages = append(messages, line)
		overallValid = false
	} else {
		// line := fmt.Sprintf("Crypticfile '%v' exists.\n", filename)
		// messages = append(messages, line)
	}

	if !publicKeyExists {
		line := fmt.Sprintf("Public key '%v' does not exist.\n", publicKey)
		messages = append(messages, line)
		overallValid = false
	} else {
		// line := fmt.Sprintf("Public key '%v' exists.\n", publicKey)
		// messages = append(messages, line)
	}

	if !privateKeyExists {
		line := fmt.Sprintf("Private key '%v' does not exist.\n", privateKey)
		messages = append(messages, line)
		overallValid = false
	} else {
		// line := fmt.Sprintf("Private key '%v' exists.\n", privateKey)
		// messages = append(messages, line)
	}

	if publicKeyExists && privateKeyExists {
		// try to encrypt/decrypt something
		plain := "Hello, World"
		encrypted := Encrypt(plain, publicKey)
		decrypted := Decrypt(encrypted, privateKey)
		if plain == decrypted {
			// line := fmt.Sprintf("Encrypt/Decrypt works.\n")
			// messages = append(messages, line)
		} else {
			line := fmt.Sprintf("Encrypt/Decrypt not working.\n")
			messages = append(messages, line)
			overallValid = false
		}

	} else {
		messages = append(messages, "\nPublic/private keys do not exist, try the following\n\n")
		line := fmt.Sprintf("    ssh-keygen -m pem -f ~/.ssh/id_rsa\n")
		messages = append(messages, line)
		line = fmt.Sprintf("    ssh-keygen -f ~/.ssh/id_rsa.pub -e -m pem > ~/.ssh/id_rsa.pem\n\n")
		messages = append(messages, line)
	}

	for _, line := range messages {
		fmt.Printf(line)
	}
	if overallValid {
		fmt.Printf("cryptic verify : OK.\n")
	} else {
		// fmt.Printf("cryptic verify: NOT OK.\n")
	}

	// fmt.Printf("%v =%v, exists=%v\n", CRYPTIC_ENCRYPTION_ENABLED, privateKeyExists)
}

func DoInfo(cli *goutils.CLI) {
	filename := goutils.GetEnvOrDefault(CRYPTIC_FILE, "~/.Crypticfile")
	pubKey := goutils.GetEnvOrDefault(CRYPTIC_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privKey := goutils.GetEnvOrDefault(CRYPTIC_PRIVATE_KEY, "~/.ssh/id_rsa")

	fmt.Printf("\n %v          : %v\n", CRYPTIC_FILE, filename)
	fmt.Printf(" %v    : %v\n", CRYPTIC_PUBLIC_KEY, pubKey)
	fmt.Printf(" %v   : %v\n\n", CRYPTIC_PRIVATE_KEY, privKey)

	fmt.Printf("\n%v\n", GLOBAL_SSH_KEYGEN_USAGE)
}

func DoGet(cli *goutils.CLI) {
	command := cli.GetCommand()
	key := cli.GetStringOrDie(command)
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
	command := cli.GetCommand()
	key := cli.GetStringOrDie(command)
	description := cli.GetStringOrDefault("-d", "")
	password := ""
	if cli.IndexOf("-value") > -1 {
		password = cli.GetStringOrDefault("-value", "")
		if password == "" {
			fmt.Printf("Error, -value cannot be empty.")
			os.Exit(1)
		}
	} else {
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password = string(bytePassword)
	}
	value := password
	db.Put(key, value, description)
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
		fmt.Printf("--------------------------------------------------------------------------\n")
		fmt.Printf("| Key       | Description          | Last Updated      | Created         |\n")
		fmt.Printf("--------------------------------------------------------------------------\n")
		for _, key := range keys {
			// fmt.Printf("%v\n", key)
			entry := data.Entries[key]
			createdStr := entry.Created.Format("January 2, 2006")
			updatedStr := entry.LastUpdated.Format("January 2, 2006")
			fmt.Printf("| %-10v| %-20v | %v | %v |\n", key, entry.Description, updatedStr, createdStr)
		}
		fmt.Printf("--------------------------------------------------------------------------\n")
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
