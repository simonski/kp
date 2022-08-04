package main

import (
	"fmt"
	"os"
	"testing"

	cli "github.com/simonski/cli"
)

func TestMain(t *testing.T) {

	command := "fooo"
	cli := cli.New(os.Args)
	if command == "test" {
		db := LoadDB()
		value := cli.GetStringOrDie(command)
		valueEnc, _ := db.Encrypt(value)
		fmt.Printf("Encrypt('%v')=\n%v\n", value, valueEnc)
		valueDec, _ := db.Decrypt(valueEnc)
		fmt.Printf("\n\nDecrypt\n '%v'\n", valueDec)
		os.Exit(0)
	}

}
