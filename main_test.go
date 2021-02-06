package main

import (
	"fmt"
	"os"
	"testing"

	goutils "github.com/simonski/goutils"
)

func main_test(t *testing.T) {
	fmt.Printf("OK")

	command := "fooo"
	cli := goutils.NewCLI(os.Args)
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
