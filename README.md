# Cryptic

A terminal tool to manage key/pairs.

## Install

Install via `go get`

	go get github.com/simonski/cryptic

This will install `cryptic` onto your `$GOBIN`. Please ensure `$GOBIN` is on your `$PATH`.

## Usage

`cryptic` stores keypairs to a `~/.Crypticfile`

The file itself is plaintext, the values of the keys are encrypted using your public key.

## Initial Setup

This will verify your environment variables, file locations and encryption keys.

	cryptic verify


## Store a key

Store a key/value

	cryptic <keyname> <value>

Retrieve the value of a to your clipboard

	cryptic <keyname>

List all keys

	cryptic ls

Remove a key

	crpytic rm <keyname>

Remove all keys

	cryptic clear

Get help

	cryptic

# Environment variables

You can optionally override settings such as encryption, location of files by setting the following environment variables:

|name|dedscription|default value|
-----|------------|-------------|
`$CRYTPIC_FILE`|The file keypairs are stored|`~/.Crypticfile`
`$CRYTPIC_PUBLICKEY`|The public key used for encryption|`~/.ssh/id_rsa.pem`
`$CRYTPIC_PRIVATEKEY`|The file keypairs are stored|`~/.ssh/id_rsa

# Initial Setup

1. Create your keypair if you haven't already

	ssh-keygen

2. Create a pem readable public key

	ssh-keygen -f ~/.ssh/id_rsa.pub -e -m pem > ~/.ssh/id_rsa.pem

3. Verify you can create a key/pair 

	cryptic verify

`
