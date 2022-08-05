# KP

A terminal tool to manage key/pairs. I use it to manage *temporary* key/pairs on *controlled* hardware.

This is OSS - if you want to contribute please read the [DEVELOPER_GUIDE.md](DEVELOPER_GUIDE.md).

## Install

Install via `go get`

```bash
go get github.com/simonski/kp
```

This will install `kp` onto your `$GOBIN`. Please ensure `$GOBIN` is on your `$PATH`.

## Setup

Once you've installed kp and can type `kp version`, you will need to configure kp.

```bash
kp verify
```

By default, `kp` stores keypairs to a `~/.kpfile`.  This can be controlled with the environment variable `KP_FILE`

## Encryption

The `~/.kpfile` itself is plaintext, the password is encrypted.

|name|purpose|default value|
|----|-------|-------------|
`KP_KEY`|Flename of private key|`~/.ssh/kp.id_rsa`

Create your encryption key

```bash
ssh-keygen -b 2048 -t rsa -N "" -m pkcs8 -f ~/.ssh/kp.id_rsa
```

Or re-use one? - up to you:

```bash
export KP_KEY=~/path/to/id_rsa
```

Finally, confirm kp is setup properly:

```bash
kp verify
```

Assuming you get a "KP is setup correctly." message, you can then use `kp` in the following manner - if you don't, it will explain what needs changing in the `verify` command itself.

## Usage

### Store a key/value

```bash
kp put <keyname>
>> STDIN value
```

### Retrieve the value of a key to your clipboard

```bash
kp get <keyname>
```

### List all keys

```bash
kp ls
```

### Search for an entry

```bash
kp search widget
```

### Update a key

```bash
kp update <key> 
 -type         "type"
 -url          "url"
 -username     "username"
 -description  "description"
 -notes        "notes"
```

### Remove a key

```bash
kp rm <keyname>
```

Get help on any command:

```bash
kp
```

## Environment variables

You can optionally override settings such as encryption, location of files by setting the following environment variables:

|name|description|default value|
-----|------------|-------------|
`$KP_FILE`|The file keypairs are stored|`~/.kpfile`
`$KP_KEY`|The encryption key|`~/.ssh/kp.id_rsa`

# Releases

I use github actions to create a crossplatform release binary on a tag.
