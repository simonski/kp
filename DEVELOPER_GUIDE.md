# Developer Guide

## Initial Setup

Install dependencies

```bash
make setup
```

## Building

To build locally

```bash
make
```

To see the help targets

```bash
make help
```

I use some workflows in github actions (./github/workflows) to compile and test.  

## PRs

Use github PRs and vefify using "Compile and Test" actions that the PR is good.

## Extending

### TermUI

`kp -g` or `KP_GUI=1` should launch the fullscreen terminal kp.  

- Q will always insta-quits
- Enter on a key always loads into memory
- E on a key goes to edit more

Consider use of password-based encryption.  This would change the flow down to requesting a password before retrieving the sensitive information on the entry.  Currently the public key encryption is a nicer UX, less invasive but does have the risk of key exposure.

### Encryption via password

Consider use of password-based encryption.  This would change the flow down to requesting a password before retrieving the sensitive information on the entry.  Currently the public key encryption is a nicer UX, less invasive but does have the risk of key exposure.

See [https://bruinsslot.jp/post/golang-crypto/]
(<https://bruinsslot.jp/post/golang-crypto/>)

```bash
kp init password
kp init key
```

### Per-key passwords

```bash
kp put key -p
<password to store>
<password to encrypt with>
```
