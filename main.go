// kp is a toy keypair manager
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"

	clipboard "github.com/atotto/clipboard"
	figure "github.com/common-nighthawk/go-figure"
	cli "github.com/simonski/cli"
	goutils "github.com/simonski/goutils"
	crypto "github.com/simonski/goutils/crypto"
	terminal "golang.org/x/term"
)

func main() {
	graphics_env := cli.GetEnvOrDefault("KP_GUI", "0") == "1"
	cli := cli.New(os.Args)
	graphics_cli := cli.IndexOf("-g") > -1
	command := cli.GetCommand()

	if graphics_cli || graphics_env {
		DoGraphics(cli)
		return
	} else if command == "help" {
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
		searchTerm := cli.GetStringOrDefault(command, "")
		if searchTerm == "-a" {
			searchTerm = cli.GetStringOrDefault(searchTerm, "")
		}
		DoList(cli, searchTerm)
	} else if isHide(command) {
		DoHide(cli)
	} else if isShow(command) {
		DoShow(cli)
	} else if isVersion(command) {
		DoVersion(cli)
	} else if isPut(command, cli) {
		DoPut(cli)
	} else if isEncrypt(command) {
		DoEncrypt(cli)
	} else if isDecrypt(command) {
		DoDecrypt(cli)
	} else if isUpdate(command) {
		DoUpdate(cli)
	} else if isRename(command) {
		DoRename(cli)
	} else if isGet(command, cli) {
		DoGet(cli)
		DoDescribe(cli)
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

func isRename(command string) bool {
	return command == "rename"
}

func isVersion(command string) bool {
	return command == "version"
}

func isList(command string) bool {
	return command == "ls" || command == "list"
}

func isHide(command string) bool {
	return command == "hide"
}

func isShow(command string) bool {
	return command == "show"
}

// A 'get' is basically not a a list, delete or a put
func isGet(command string, c *cli.CLI) bool {
	return command == "get"
}

func isPut(command string, c *cli.CLI) bool {
	return command == "put"
}

func DoGraphics(c *cli.CLI) {
	filename := cli.GetEnvOrDefault(KP_FILE, DEFAULT_DB_FILE)
	privKey := cli.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)
	db := NewKPDB(filename, privKey)
	gui := NewGUI(db)
	gui.Run()
}

func LoadDB() *KPDB {
	filename := cli.GetEnvOrDefault(KP_FILE, DEFAULT_DB_FILE)
	privKey := cli.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)
	db := NewKPDB(filename, privKey)
	return db
}

func DoEncrypt(c *cli.CLI) {
	privateKey := cli.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)
	command := c.GetCommand()
	value := c.GetStringOrDie(command)
	result, err := crypto.Encrypt(value, privateKey)
	if err != nil {
		fmt.Printf("Problem decrypting:\n%v\n", err)
		os.Exit(1)
	} else {
		fmt.Println(result)
	}
}

func DoDecrypt(c *cli.CLI) {
	privateKey := cli.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)
	command := c.GetCommand()
	value := c.GetStringOrDie(command)
	result, _ := crypto.Decrypt(value, privateKey)
	fmt.Println(result)
}

// DoVerify performs verification of ~/.KPfile, encryption/decryption using
// specified keys
func DoVerify(c *cli.CLI, printFailuresToStdOut bool) bool {
	overallValid := true
	kpFilename := cli.GetEnvOrDefault(KP_FILE, DEFAULT_DB_FILE)
	privateKeyFilename := cli.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)

	filenameExists := goutils.FileExists(goutils.EvaluateFilename(kpFilename))
	privateKeyExists := goutils.FileExists(goutils.EvaluateFilename(privateKeyFilename))

	messages := make([]string, 0)
	messages = append(messages, fmt.Sprintf("%v   : %v, exists=%v\n", KP_FILE, kpFilename, filenameExists))
	messages = append(messages, fmt.Sprintf("%v    : %v, exists=%v\n", KP_KEY, privateKeyFilename, privateKeyExists))

	if !privateKeyExists {
		overallValid = false
		messages = append(messages, "\nEncryption key does not exist, try the following\n\n")
		line := fmt.Sprintf("    %v\n\n", GetSSHCommand(privateKeyFilename))
		messages = append(messages, line)
	}

	if overallValid {
		// try to encrypt/decrypt something
		plain := "Hello, World"
		encrypted, err := crypto.Encrypt(plain, privateKeyFilename)
		if err != nil {
			line := fmt.Sprintf("Error encrypting: %v\n", err)
			messages = append(messages, line)
			overallValid = false
		}
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

func DoInfo(c *cli.CLI) {

	filename := cli.GetEnvOrDefault(KP_FILE, DEFAULT_DB_FILE)
	privKey := cli.GetEnvOrDefault(KP_KEY, DEFAULT_KEY_FILE)

	fmt.Printf("\nKP is currently using the following values:\n")
	fmt.Printf("\n%v  : %v\n", KP_FILE, filename)
	fmt.Printf("%v   : %v\n", KP_KEY, privKey)
	msg := strings.ReplaceAll(GLOBAL_SSH_KEYGEN_USAGE, "TOKEN_DEFAULT_DB_FILE", filename)
	msg = strings.ReplaceAll(msg, "TOKEN_DEFAULT_KEY_FILE", privKey)
	msg = strings.ReplaceAll(msg, "TOKEN_DEFAULT_SSH_COMMAND", GetSSHCommand(privKey))
	fmt.Printf("\n%v\n", msg)

}

func DoGet(c *cli.CLI) {
	command := c.GetCommand()
	key := c.GetStringOrDie(command)
	db := LoadDB()
	entry, exists := db.GetDecrypted(key)
	if exists {
		value := entry.Value
		err := clipboard.WriteAll(value)
		if err != nil {
			fmt.Printf("%v", err)
		}
		if c.IndexOf("-stdout") > -1 {
			fmt.Printf("%v\n", value)
		}
	} else {
		fmt.Printf("'%v' does not exist.\n", key)
		os.Exit(1)
	}
}

func DoPut(c *cli.CLI) {
	db := LoadDB()
	command := c.GetCommand()
	key := c.GetStringOrDie(command)
	if len(key) > 125 {
		fmt.Printf("Error, key must be <= 25 characters.\n")
		os.Exit(1)
	} else if len(key) == 0 {
		fmt.Printf("Error, key cannot be empty.")
		os.Exit(1)
	}
	entry, _ := db.GetDecrypted(key)
	entry.Key = key
	entry.Description = c.GetStringOrDefault("-description", entry.Description)
	entry.Type = c.GetStringOrDefault("-type", entry.Type)
	entry.Notes = c.GetStringOrDefault("-note", entry.Notes)
	entry.Url = c.GetStringOrDefault("-url", entry.Url)
	entry.Username = c.GetStringOrDefault("-username", entry.Username)

	password := ""
	if c.IndexOf("-value") > -1 {
		password = c.GetStringOrDefault("-value", "")
		if password == "" {
			fmt.Printf("Error, -value cannot be empty.")
			os.Exit(1)
		}
	} else {
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password = string(bytePassword)
	}
	if password != "" {
		entry.Value = password
	}
	db.Put(entry)
	db.Save()
}

func DoUpdate(c *cli.CLI) {
	db := LoadDB()
	command := c.GetCommand()
	key := c.GetStringOrDie(command)
	entry, _ := db.GetDecrypted(key)
	entry.Description = c.GetStringOrDefault("-description", entry.Description)
	entry.Type = c.GetStringOrDefault("-type", entry.Type)
	entry.Notes = c.GetStringOrDefault("-note", entry.Notes)
	entry.Url = c.GetStringOrDefault("-url", entry.Url)
	entry.Username = c.GetStringOrDefault("-username", entry.Username)
	db.Put(entry)
	db.Save()
}

func DoHide(c *cli.CLI) {
	db := LoadDB()
	command := c.GetCommand()
	key := c.GetStringOrDie(command)
	entry, _ := db.GetDecrypted(key)
	entry.Hidden = true
	db.Put(entry)
	db.Save()
}

func DoShow(c *cli.CLI) {
	db := LoadDB()
	command := c.GetCommand()
	key := c.GetStringOrDie(command)
	entry, _ := db.GetDecrypted(key)
	entry.Hidden = false
	db.Put(entry)
	db.Save()
}

func DoRename(c *cli.CLI) {
	db := LoadDB()
	command := c.GetCommand()
	old_key := c.GetStringOrDie(command)
	entry, exists := db.GetDecrypted(old_key)
	if !exists {
		fmt.Printf("No such entry '%v'\n", old_key)
		os.Exit(1)
	}

	new_key := c.GetStringOrDie(old_key)
	_, other_exists := db.GetDecrypted(new_key)
	if other_exists {
		fmt.Printf("Entry '%v' already exists.\n", new_key)
		os.Exit(1)
	}
	entry.Key = new_key
	db.Delete(old_key)
	db.Put(entry)
	db.Save()
}

func DoDescribe(c *cli.CLI) {
	db := LoadDB()
	command := c.GetCommand()
	key := c.GetStringOrDie(command)
	entry, exists := db.GetDecrypted(key)
	if exists {
		fmt.Printf("Key          : %v\n", entry.Key)
		fmt.Printf("Description  : %v\n", entry.Description)
		fmt.Printf("Username     : %v\n", entry.Username)
		fmt.Printf("Url          : %v\n", entry.Url)
		fmt.Printf("Created      : %v\n", entry.Created.Format(time.RFC822))
		fmt.Printf("Last Updated : %v\n", entry.LastUpdated.Format(time.RFC822))
		fmt.Printf("Type         : %v\n", entry.Type)
		fmt.Printf("Notes        : %v\n", entry.Notes)
	} else {
		fmt.Printf("Error, no entry '%v'\n", entry.Key)
	}

}

func DoList(c *cli.CLI, searchTerm string) {
	db := LoadDB()
	data := db.GetData()
	includeHidden := c.IndexOf("-a") > -1

	if len(data.Entries) == 0 {
		fmt.Printf("DB is empty.\n")
	} else {
		maxLength := 0
		for key := range data.Entries {
			entry := data.Entries[key]
			if !includeHidden && entry.Hidden {
				continue
			}
			maxLength = goutils.Max(len(key), maxLength)
		}

		keys := make([]string, 0)
		max_key := len("Key") + 1
		max_url := len("Url") + 1
		max_username := len("Username") + 1
		max_type := len("Type") + 1
		max_description := len("Description") + 1
		max_notes := len("Notes") + 1

		for key, entry := range data.Entries {
			if !includeHidden && entry.Hidden {
				continue
			}
			keys = append(keys, key)
			max_key = goutils.Max(len(entry.Key)+1, max_key)
			max_url = goutils.Max(len(entry.Url)+1, max_url)
			max_username = goutils.Max(len(entry.Username)+1, max_username)
			max_type = goutils.Max(len(entry.Type)+1, max_type)
			max_description = goutils.Max(len(entry.Description)+1, max_description)
			// max_notes = goutils.Max(len(entry.Notes)+1, max_notes)
		}
		sort.Strings(keys)

		date_width := len(time.Now().Format(time.RFC822)) + 1
		width := max_key + max_url + max_username + max_type + max_description + max_notes + (2 * date_width) + 8

		line := strings.Repeat("-", width+1)
		fmt.Println(line)

		entry_line := fmt.Sprintf("|%v|%v|%v|%v|%v|%v|%v|%v|",
			goutils.RPadToFixedLength("Key", " ", max_key),
			goutils.RPadToFixedLength("Username", " ", max_username),
			goutils.RPadToFixedLength("Url", " ", max_url),
			goutils.RPadToFixedLength("Type", " ", max_type),
			goutils.RPadToFixedLength("Description", " ", max_description),
			goutils.RPadToFixedLength("Notes", " ", max_notes),
			goutils.RPadToFixedLength("Updated", " ", date_width),
			goutils.RPadToFixedLength("Created", " ", date_width))
		fmt.Println(entry_line)
		fmt.Println(line)

		for _, entry := range db.GetEntriesSortedByUpdatedThenKey() {
			if !includeHidden && entry.Hidden {
				continue
			}

			if searchTerm != "" {
				found := strings.Contains(entry.Key, searchTerm)
				found = found || strings.Contains(entry.Description, searchTerm)
				found = found || strings.Contains(entry.Notes, searchTerm)
				found = found || strings.Contains(entry.Type, searchTerm)
				found = found || strings.Contains(entry.Username, searchTerm)
				found = found || strings.Contains(entry.Url, searchTerm)
				if !found {
					continue
				}
			}
			notes := entry.Notes
			if notes != "" {
				notes = "***"
			}
			created := entry.Created.Format(time.RFC822)
			updated := entry.LastUpdated.Format(time.RFC822)
			desc := entry.Description

			if len(desc) > max_description {
				desc = desc[0:max_description-3] + "..."
			}

			entry_line := fmt.Sprintf("|%v|%v|%v|%v|%v|%v|%v|%v",
				goutils.RPadToFixedLength(entry.Key, " ", max_key),
				goutils.RPadToFixedLength(entry.Username, " ", max_username),
				goutils.RPadToFixedLength(entry.Url, " ", max_url),
				goutils.RPadToFixedLength(entry.Type, " ", max_type),
				goutils.RPadToFixedLength(desc, " ", max_description),
				goutils.RPadToFixedLength(notes, " ", max_notes),
				goutils.RPadToFixedLength(created, " ", date_width),
				goutils.RPadToFixedLength(updated, " ", date_width))
			fmt.Println(entry_line)

		}
		fmt.Println(line)
	}
}

func DoDelete(c *cli.CLI) {
	command := c.GetCommand()
	key := c.GetStringOrDefault(command, "")
	if key == "" {
		USAGE := "kp rm [key]"
		fmt.Printf("%v\n", USAGE)
	}
	db := LoadDB()
	_, exists := db.GetDecrypted(key)
	if !exists {
		fmt.Printf("Error, '%v' does not exist.\n", key)
		os.Exit(1)
	}
	db.Delete(key)
	db.Save()

}

func DoUsage(c *cli.CLI) {
	fmt.Print(GLOBAL_USAGE)
}

func DoVersion(c *cli.CLI) {
	fmt.Printf("%v\n", VERSION)
}
