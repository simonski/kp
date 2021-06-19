package main

import (
	"bytes"
	// "fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Terminal struct {
}

func NewTerminal() *Terminal {
	return &Terminal{}
}

func (*Terminal) Width() int {
	// fmt.Printf("Environ: %v\n\n", os.Environ())

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	v := strings.Replace(out.String(), "\n", "", -1)
	splits := strings.Split(v, " ")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("tput output is %v\n", out.String())

	ival, err := strconv.Atoi(splits[1])
	if err != nil {
		panic(err)
	}
	return ival
}
