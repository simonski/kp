# KP

A terminal tool to manage encrypted key/value pairs.

## Install

### macOS (Homebrew)

```bash
brew install simonski/tap/kp
```

### From source

```bash
go install github.com/simonski/kp@latest
```

Ensure `$GOBIN` is on your `$PATH`.

## Quick Start

After installing, kp needs an RSA key for encryption. Generate one:

```bash
ssh-keygen -b 2048 -t rsa -N "" -m pkcs8 -f ~/.ssh/kp.id_rsa
```

Verify everything is working:

```bash
kp verify
```

You should see "KP is setup correctly." See [QUICKSTART.md](QUICKSTART.md) for a full walkthrough.

## Usage

```bash
kp put mykey                  # store a value (prompts for input)
kp put mykey -value "secret"  # store a value inline
kp put mykey -random 32       # store a generated 32-char password
kp get mykey                  # copy value to clipboard
kp get mykey -stdout          # print value to stdout
kp ls                         # list all keys
kp ls -a                      # list all keys (including hidden)
kp ls widget                  # search for keys matching "widget"
kp rm mykey                   # delete a key
kp rename old new             # rename a key
kp hide mykey                 # hide a key from default listing
kp show mykey                 # unhide a key
kp open mykey                 # open the URL associated with a key
kp info                       # show current configuration
kp version                    # print version
```

### Updating metadata

```bash
kp update mykey -url "https://example.com" -username "me" -description "My account" -notes "some notes" -type "login"
```

### Tags

```bash
kp tag mykey work             # add a tag
kp untag mykey work           # remove a tag
```

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `KP_FILE` | `~/.kpfile` | Path to the encrypted key/pair database |
| `KP_KEY` | `~/.ssh/kp.id_rsa` | Path to RSA private key for encryption |
| `KP_GUI` | `0` | Set to `1` to launch TUI mode |

## TUI Mode

```bash
kp -g
```

Or set `KP_GUI=1`. Q quits, Enter copies value to clipboard, E enters edit mode.

## Development

See [DEVELOPER_GUIDE.md](DEVELOPER_GUIDE.md) for contribution guidelines.

Tickets are managed with [tk](https://github.com/simonski/ticket) (`brew install simonski/tap/ticket`). Run `tk list` to see open work.

## License

MIT
