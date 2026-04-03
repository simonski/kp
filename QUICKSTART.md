# Quick Start

This guide walks you through setting up kp after installation.

## 1. Install kp

### macOS (Homebrew)

```bash
brew install simonski/tap/kp
```

### From source

```bash
go install github.com/simonski/kp@latest
```

## 2. Create an encryption key

kp encrypts all stored values using an RSA key. Generate one:

```bash
ssh-keygen -b 2048 -t rsa -N "" -m pkcs8 -f ~/.ssh/kp.id_rsa
```

This creates `~/.ssh/kp.id_rsa` (private key) and `~/.ssh/kp.id_rsa.pub` (public key). kp uses the private key for both encryption and decryption.

If you want to use a different key path:

```bash
export KP_KEY=~/.ssh/my_other_key
```

## 3. Verify the setup

```bash
kp verify
```

This checks that:
- The encryption key exists at `~/.ssh/kp.id_rsa` (or `$KP_KEY`)
- Encrypt/decrypt round-trips work correctly

You should see:

```
KP_FILE   : ~/.kpfile, exists=false
KP_KEY    : ~/.ssh/kp.id_rsa, exists=true
KP is setup correctly.
```

The database file (`~/.kpfile`) is created automatically the first time you store a value. It's fine for `exists=false` at this point.

If you see **"Failed to verify encryption."** it means the encryption key is missing. Run the `ssh-keygen` command from step 2.

## 4. Store your first value

```bash
kp put my-email
```

You'll be prompted to enter a value (input is hidden). Or store it inline:

```bash
kp put my-email -value "me@example.com"
```

Or generate a random password:

```bash
kp put my-password -random 32
```

## 5. Retrieve a value

Copy to clipboard:

```bash
kp get my-password
```

Print to stdout:

```bash
kp get my-password -stdout
```

## 6. List your keys

```bash
kp ls
```

## Troubleshooting

### "Failed to verify encryption."

The RSA key doesn't exist at the expected path. Either:

1. Generate one: `ssh-keygen -b 2048 -t rsa -N "" -m pkcs8 -f ~/.ssh/kp.id_rsa`
2. Or point kp at an existing key: `export KP_KEY=/path/to/your/key`

Then run `kp verify` to confirm.

### Key format requirements

The key must be RSA in PKCS#8 format. Keys generated with the `-m pkcs8` flag work correctly. If you have an existing RSA key in a different format, generate a new one with the command above.
