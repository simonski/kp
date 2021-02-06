# Cryptic

A terminal tool to manage key/pairs. I use it to manage *temporary* passwords on *controlled* hardware.

## Install

Install via `go get`

	go get github.com/simonski/cryptic

This will install `cryptic` onto your `$GOBIN`. Please ensure `$GOBIN` is on your `$PATH`.

## Usage

`cryptic` stores keypairs to a `~/.Crypticfile`

The file itself is plaintext, the values of the keys are encrypted using your public key.


## Verify Setup

(Optional) this step will verify your installation, environment variables, file locations and encryption keys. You don't *have* to do this as the defaults *should* work on modern mac, Windows and Linux variants.

	cryptic verify

You may need to create your encryption keys

	ssh-keygen

Create a pem readable public key

	ssh-keygen -f ~/.ssh/id_rsa.pub -e -m pem > ~/.ssh/id_rsa.pem

Finally, confirm cryptic is setup properly:

	cryptic verify

Assuming you get a "Verification Success" message, you can then use `cryptic` in the following manner - if you don't, it will explain what needs changing in the `verify` command itself.

## Store a key

Store a key/value

	cryptic <keyname> <value>

Retrieve the value of a to your clipboard

	cryptic <keyname>

> Note: storing and retrieving keys uses the key name as the command (no "get" or "set").  When you retrieve a value, it won't write to STDOUT - your clipboard will contain the value.

List all keys

	cryptic ls

Remove a key

	crpytic rm <keyname>

Remove all keys

	cryptic clear

Get help on any command:

	cryptic

# Environment variables

You can optionally override settings such as encryption, location of files by setting the following environment variables:

|name|dedscription|default value|
-----|------------|-------------|
`$CRYPTIC_FILE`|The file keypairs are stored|`~/.Crypticfile`
`$CRYPTIC_PUBLICKEY`|The public key used for encryption|`~/.ssh/id_rsa.pem`
`$CRYPTIC_PRIVATEKEY`|The file keypairs are stored|`~/.ssh/id_rsa`