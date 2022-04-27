package main

import (
	"fmt"
	"os"
	"testing"

	cli "github.com/simonski/cli"
)

func main_test(t *testing.T) {

	command := "fooo"
	cli := cli.New(os.Args)
	if command == "test" {
		db := LoadDB()
		value := cli.GetStringOrDie(command)
		valueEnc := db.Encrypt(value)
		fmt.Printf("Encrypt('%v')=\n%v\n", value, valueEnc)
		valueDec := db.Decrypt(valueEnc)
		fmt.Printf("\n\nDecrypt\n '%v'\n", valueDec)
		os.Exit(0)
	}

}
