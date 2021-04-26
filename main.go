package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
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
		DoVerify(cli, true)
	} else if isVersion(command) {
		DoVersion(cli)
	} else if isClear(command) {
		DoClear(cli)
	} else if isList(command) {
		DoList(cli)
	} else if isPut(command, cli) {
		ok := DoVerify(cli, false)
		if ok {
			DoPut(cli)
		}
	} else if isGet(command, cli) {
		ok := DoVerify(cli, false)
		if ok {
			DoGet(cli)
		}
	} else if isDelete(command) {
		DoDelete(cli)
	} else if command != "" {
		fmt.Printf("kp %v: unknown command\n", command)
		fmt.Printf("Run 'kp help' for usage.\n")
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

func LoadDB() *KPDB {
	filename := goutils.GetEnvOrDefault(KP_FILE, "~/.KPfile")
	pubKey := goutils.GetEnvOrDefault(KP_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privKey := goutils.GetEnvOrDefault(KP_PRIVATE_KEY, "~/.ssh/id_rsa")
	encryptionEnabled := goutils.GetEnvOrDefault(KP_ENCRYPTION_ENABLED, "1") == "1"
	db := NewKPDB(filename, pubKey, privKey, encryptionEnabled)
	return db
}

// DoVerify performs verification of ~/.KPfile, encryption/decryption using
// specified keys
func DoVerify(cli *goutils.CLI, printFailuresToStdOut bool) bool {
	overallValid := true
	filename := goutils.GetEnvOrDefault(KP_FILE, "~/.KPfile")
	publicKey := goutils.GetEnvOrDefault(KP_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privateKey := goutils.GetEnvOrDefault(KP_PRIVATE_KEY, "~/.ssh/id_rsa")

	// filenameExists := goutils.FileExists(goutils.EvaluateFilename(filename))
	publicKeyExists := goutils.FileExists(goutils.EvaluateFilename(publicKey))
	privateKeyExists := goutils.FileExists(goutils.EvaluateFilename(privateKey))

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("%v          %v\n", KP_FILE, filename))
	messages = append(messages, fmt.Sprintf("%v    %v\n", KP_PUBLIC_KEY, publicKey))
	messages = append(messages, fmt.Sprintf("%v   %v\n", KP_PRIVATE_KEY, privateKey))

	// overallValid = filenameExists && publicKeyExists && privateKeyExists
	// if !filenameExists {
	// 	line := fmt.Sprintf("KPfile '%v' does not exist.\n", filename)
	// 	messages = append(messages, line)
	// 	overallValid = false
	// } else {
	// 	// line := fmt.Sprintf("KPfile '%v' exists.\n", filename)
	// 	// messages = append(messages, line)
	// }

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

	if printFailuresToStdOut {
		for _, line := range messages {
			fmt.Printf(line)
		}
	}

	if overallValid {
		if printFailuresToStdOut {
			fmt.Printf("kp verify : OK.\n")
		}
	} else {
		// fmt.Printf("kp verify: NOT OK.\n")
	}

	return overallValid
	// fmt.Printf("%v =%v, exists=%v\n", KP_ENCRYPTION_ENABLED, privateKeyExists)
}

func DoInfo(cli *goutils.CLI) {
	filename := goutils.GetEnvOrDefault(KP_FILE, "~/.KPfile")
	pubKey := goutils.GetEnvOrDefault(KP_PUBLIC_KEY, "~/.ssh/id_rsa.pem")
	privKey := goutils.GetEnvOrDefault(KP_PRIVATE_KEY, "~/.ssh/id_rsa")

	fmt.Printf("\nKP is currently using the following values:\n")
	fmt.Printf("\n %v          : %v\n", KP_FILE, filename)
	fmt.Printf(" %v    : %v\n", KP_PUBLIC_KEY, pubKey)
	fmt.Printf(" %v   : %v\n\n", KP_PRIVATE_KEY, privKey)

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
	if len(key) > 25 {
		fmt.Printf("Error, key must be <= 25 characters.\n")
		os.Exit(1)
	}
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

		key_width := 50
		desc_width := 50
		date_width := 25
		width := key_width + desc_width + (2 * date_width) + 4

		line := strings.Repeat("-", width) + "\n"
		fmt.Printf(line)
		header := fmt.Sprintf("| Key%v| Description%v| Updated%v| Created%v|\n", strings.Repeat(" ", key_width-len("Key")-1), strings.Repeat(" ", desc_width-len("Description")-1), strings.Repeat(" ", date_width-len("Updated")-1), strings.Repeat(" ", date_width-len("Created")-1))
		fmt.Printf(header)
		fmt.Printf(line)
		max_key_length := key_width - 5
		max_description_length := desc_width - 5
		for _, key := range keys {
			entry := data.Entries[key]
			created := entry.Created.Format("January 2, 2006")
			updated := entry.LastUpdated.Format("January 2, 2006")
			desc := entry.Description

			if len(key) > max_key_length {
				key = key[0:max_key_length] + "..."
			}

			if desc == "" {
				desc = "No description."
			}
			if len(desc) > max_description_length {
				desc = desc[0:max_description_length] + "..."
			}

			keyExtra := key_width - 1 - len(key)
			descExtra := desc_width - 1 - len(desc)
			updatedExtra := date_width - 1 - len(updated)
			createdExtra := date_width - 1 - len(created)

			if keyExtra < 0 {
				keyExtra = 0
			}
			if descExtra < 0 {
				descExtra = 0
			}
			if updatedExtra < 0 {
				updatedExtra = 0
			}
			if createdExtra < 0 {
				createdExtra = 0
			}

			entry_line := fmt.Sprintf("| %v%v| %v%v| %v%v| %v%v|\n", key, strings.Repeat(" ", keyExtra), desc, strings.Repeat(" ", descExtra), updated, strings.Repeat(" ", updatedExtra), created, strings.Repeat(" ", createdExtra))
			fmt.Printf(entry_line)
		}
		fmt.Printf(line)
	}
}

func DoDelete(cli *goutils.CLI) {
	command := cli.GetCommand()
	key := cli.GetStringOrDefault(command, "")
	if key == "" {
		USAGE := "kp rm [key]"
		fmt.Printf("%v\n", USAGE)
	}
	db := LoadDB()
	_, exists := db.Get(key)
	if !exists {
		fmt.Printf("Error, '%v' does not exist.\n", key)
		os.Exit(1)
	}
	db.Delete(key)
	db.Save()

}

func DoUsage(cli *goutils.CLI) {
	fmt.Printf(GLOBAL_USAGE)
}

func DoVersion(cli *goutils.CLI) {
	fmt.Printf("%v\n", VERSION)
}
