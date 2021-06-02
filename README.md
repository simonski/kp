# KP

A terminal tool to manage key/pairs. I use it to manage *temporary* key/pairs on *controlled* hardware.

## Install

Install via `go get`

	go get github.com/simonski/kp

This will install `kp` onto your `$GOBIN`. Please ensure `$GOBIN` is on your `$PATH`.

## Setup

Once you've installed kp and can type `kp version`, you will need to configure kp.

	kp verify

By default, `kp` stores keypairs to a `~/.kpfile`

The file itself is plaintext, the values of the keys are *encrypted* if `KP_ENCRYTION=1`

If you want to encrypt your data, create your encryption keys

	ssh-keygen

Create a pem readable public key

	ssh-keygen -f ~/.ssh/id_rsa.pub -e -m pem > ~/.ssh/id_rsa.pem

Finally, confirm kp is setup properly:

	kp verify

Assuming you get a "Verification Success" message, you can then use `kp` in the following manner - if you don't, it will explain what needs changing in the `verify` command itself.

## Store a key

Store a key/value

	kp put <keyname> [-m message]
	>> STDIN value

Retrieve the value of a key to your clipboard

	kp get <keyname>

List all keys

	kp ls

Remove a key

	crpytic rm <keyname>

Remove all keys

	kp clear

Get help on any command:

	kp

# Environment variables

You can optionally override settings such as encryption, location of files by setting the following environment variables:

|name|dedscription|default value|
-----|------------|-------------|
`$KP_ENCRYPTION`|1 or 0, indicates if encrytion is used.|`~/.KPfile`
`$KP_FILE`|The file keypairs are stored|`~/.kpfile`
`$KP_PUBLICKEY`|The public key used for encryption|`~/.ssh/id_rsa.pem`
`$KP_PRIVATEKEY`|The file keypairs are stored|`~/.ssh/id_rsa`

# Releases

I use github actions to create a crossplatform release binary on a tag.
