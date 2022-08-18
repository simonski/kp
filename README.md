# KP

A terminal tool to manage key/pairs.

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

### List all keys (including hidden)

```bash
kp ls -a
```

### Search for an entry

```bash
kp ls widget
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

### Hide a key

```bash
kp hide <keyname>
```

### Show the key

```bash
kp show <keyname>
```

## Environment variables

kp tries to run with sensible defaults. You can override them using the following envinronment variables:

|name|description|default value|
-----|------------|-------------|
`$KP_FILE`|The file keypairs are stored|`~/.kpfile`
`$KP_KEY`|The encryption key|`~/.ssh/kp.id_rsa`
`$KP_GUI`|Run in graphics mode (`0` or `1`)|`0`

## Releases

I use github actions to create a crossplatform release binary on a tag.
