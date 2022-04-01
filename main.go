package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"

	clipboard "github.com/atotto/clipboard"
	figure "github.com/common-nighthawk/go-figure"
	goutils "github.com/simonski/goutils"
	crypto "github.com/simonski/goutils/crypto"
	terminal "golang.org/x/term"
)

func main() {
	cli := goutils.NewCLI(os.Args)
	command := cli.GetCommand()
	if command == "help" {
		DoLogo()
		DoUsage(cli)
		return
	} else if isInfo(command) {
		DoInfo(cli)
		return
	} else if isVerify(command) {
		result := DoVerify(cli, true)
		if !result {
			fmt.Println("KP is NOT setup correctly - failed to verify encryption.")
			os.Exit(1)
		} else {
			fmt.Println("KP is setup correctly.")
			return
		}
	}

	result := DoVerify(cli, false)
	if !result {
		fmt.Println("Failed to verify encryption.")
		os.Exit(1)
	}

	if isList(command) {
		DoList(cli)
	} else if isVersion(command) {
		DoVersion(cli)
	} else if isDescribe(command) {
		DoDescribe(cli)
	} else if isPut(command, cli) {
		DoPut(cli)
	} else if isEncrypt(command) {
		DoEncrypt(cli)
	} else if isDecrypt(command) {
		DoDecrypt(cli)
	} else if isUpdate(command) {
		DoUpdateDescription(cli)
	} else if isGet(command, cli) {
		DoGet(cli)
	} else if isDelete(command) {
		DoDelete(cli)
	} else if command != "" {
		fmt.Printf("kp %v: unknown command\n", command)
		fmt.Printf("Run 'kp help' for usage.\n")
		os.Exit(1)
	} else {
		DoLogo()
		DoUsage(cli)
	}
}

func isVerify(command string) bool {
	return command == "verify"
}

func isEncrypt(command string) bool {
	return command == "encrypt"
}

func isDecrypt(command string) bool {
	return command == "decrypt"
}

func isInfo(command string) bool {
	return command == "info"
}

func isDelete(command string) bool {
	return command == "rm"
}

func isUpdate(command string) bool {
	return command == "update"
}

func isVersion(command string) bool {
	return command == "version"
}

func isList(command string) bool {
	return command == "ls"
}

func isDescribe(command string) bool {
	return command == "describe"
}

// A 'get' is basically not a a list, delete or a put
func isGet(command string, cli *goutils.CLI) bool {
	return command == "get"
}

func isPut(command string, cli *goutils.CLI) bool {
	return command == "put"
}

func LoadDB() *KPDB {
	filename := goutils.GetEnvOrDefault(KP_FILE, DEFAULT_DB_FILE)
	privKey := goutils.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)
	db := NewKPDB(filename, privKey)
	return db
}

func DoEncrypt(cli *goutils.CLI) {
	privateKey := goutils.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)
	command := cli.GetCommand()
	value := cli.GetStringOrDie(command)
	result, _ := crypto.Encrypt(value, privateKey)
	fmt.Println(result)
}

func DoDecrypt(cli *goutils.CLI) {
	privateKey := goutils.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)
	command := cli.GetCommand()
	value := cli.GetStringOrDie(command)
	result, _ := crypto.Decrypt(value, privateKey)
	fmt.Println(result)
}

// DoVerify performs verification of ~/.KPfile, encryption/decryption using
// specified keys
func DoVerify(cli *goutils.CLI, printFailuresToStdOut bool) bool {
	overallValid := true
	kpFilename := goutils.GetEnvOrDefault(KP_FILE, DEFAULT_DB_FILE)
	privateKeyFilename := goutils.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)

	filenameExists := goutils.FileExists(goutils.EvaluateFilename(kpFilename))
	privateKeyExists := goutils.FileExists(goutils.EvaluateFilename(privateKeyFilename))

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("%v   : %v, exists=%v\n", KP_FILE, kpFilename, filenameExists))
	messages = append(messages, fmt.Sprintf("%v    : %v, exists=%v\n", KP_KEY, privateKeyFilename, privateKeyExists))

	// if !filenameExists {
	// 	// line := fmt.Sprintf("KP_FILE '%v' does not exist.\n", kpFilename)
	// 	// messages = append(messages, line)
	// 	overallValid = false
	// } else {
	// 	// fmt.Printf("KP_FILE '%v' exists.\n", kpFilename)
	// }
	if !privateKeyExists {
		// line := fmt.Sprintf("KP_KEY '%v' does not exist.\n", privateKeyFilename)
		// messages = append(messages, line)
		overallValid = false
		messages = append(messages, "\nEncryption key does not exist, try the following\n\n")
		line := fmt.Sprintf("    %v", GetSSHCommand(privateKeyFilename))
		messages = append(messages, line)
	} else {
		// fmt.Printf("KP_KEY '%v' exists.\n", privateKeyFilename)

	}

	if overallValid {
		// try to encrypt/decrypt something
		plain := "Hello, World"
		encrypted, _ := crypto.Encrypt(plain, privateKeyFilename)
		decrypted, _ := crypto.Decrypt(encrypted, privateKeyFilename)
		if plain != decrypted {
			line := "Encrypt/Decrypt not working.\n"
			messages = append(messages, line)
			overallValid = false
		}

	}

	if printFailuresToStdOut {
		for _, line := range messages {
			fmt.Print(line)
		}
	}

	return overallValid
}

func DoLogo() {
	f := figure.NewColorFigure("kp", "", "blue", true)
	f.Print()
}

func DoInfo(cli *goutils.CLI) {

	filename := goutils.GetEnvOrDefault(KP_FILE, DEFAULT_DB_FILE)
	privKey := goutils.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)

	fmt.Printf("\nKP is currently using the following values:\n")
	fmt.Printf("\n%v  : %v\n", KP_FILE, filename)
	fmt.Printf("%v   : %v\n", KP_KEY, privKey)
	msg := strings.ReplaceAll(GLOBAL_SSH_KEYGEN_USAGE, "TOKEN_DEFAULT_DB_FILE", filename)
	msg = strings.ReplaceAll(msg, "TOKEN_DEFAULT_KEY_FILE", privKey)
	msg = strings.ReplaceAll(msg, "TOKEN_DEFAULT_SSH_COMMAND", GetSSHCommand(privKey))
	fmt.Printf("\n%v\n", msg)

	// t := NewTerminal()

	// sysInfo := goutils.NewSysInfo()

	// fmt.Printf("RAM         : %v\n", sysInfo.RAM)
	// fmt.Printf("CPU         : %v\n", sysInfo.CPU)
	// fmt.Printf("Cores	    : %v\n", runtime.NumCPU())
	// fmt.Printf("Disk        : %v\n", sysInfo.Disk)
	// fmt.Printf("Hostname    : %v\n", sysInfo.Hostname)
	// fmt.Printf("GOOS        : %v\n", runtime.GOOS)
	// fmt.Printf("GOARCH      : %v\n", runtime.GOARCH)
	// fmt.Printf("GOMAXPROC   : %v\n", runtime.GOMAXPROCS)
	// fmt.Printf("Columns     : %v\n", t.Width())
	// fmt.Printf("IsMacOS     : %v\n", sysInfo.IsMacOS())
	// fmt.Printf("IsLinux     : %v\n", sysInfo.IsLinux())
	// fmt.Printf("IsWindows   : %v\n", sysInfo.IsWindows())

}

func DoGet(cli *goutils.CLI) {
	command := cli.GetCommand()
	key := cli.GetStringOrDie(command)
	db := LoadDB()
	entry, exists := db.Get(key)
	if exists {
		value := entry.Value
		err := clipboard.WriteAll(value)
		if err != nil {
			fmt.Printf("%v", err)
		}
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
	} else if len(key) == 0 {
		fmt.Printf("Error, key cannot be empty.")
		os.Exit(1)
	}
	entry, exists := db.Get(key)
	var description string
	if exists {
		description = entry.Description
	} else {
		description = ""
	}
	description = cli.GetStringOrDefault("-d", description)
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

func DoUpdateDescription(cli *goutils.CLI) {
	db := LoadDB()
	command := cli.GetCommand()
	key := cli.GetStringOrDie(command)
	description := cli.GetStringOrDie(key)
	db.UpdateDescription(key, description)
	db.Save()
}

func DoDescribe(cli *goutils.CLI) {
	db := LoadDB()
	command := cli.GetCommand()
	key := cli.GetStringOrDie(command)
	description := cli.GetStringOrDie(key)
	value, _ := db.Get(key)
	db.Put(key, value.Value, description)
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
