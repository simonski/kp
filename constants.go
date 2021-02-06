package main

// VERSION is the number of this beast
const VERSION = "1.0.2-beta"

// CRYPTIC_FILE the key for the env var pointint to the file we load/save
const CRYPTIC_FILE = "CRYPTIC_FILE"

// CRYPTIC_ENCRYPTION_ENABLED the key to the public key
const CRYPTIC_ENCRYPTION_ENABLED = "CRYPTIC_ENCRYPTION_ENABLED"

// CRYPTIC_PUBLIC_KEY the key to the public key
const CRYPTIC_PUBLIC_KEY = "CRYPTIC_PUBLIC_KEY"

// CRYPTIC_PUBLIC_KEY the key to the private key
const CRYPTIC_PRIVATE_KEY = "CRYPTIC_PRIVATE_KEY"

// GLOBAL_USAGE - well, it tells me what to type
const GLOBAL_USAGE = `cryptic is a tool for using key/pairs.

Usage:

    cryptic <key | command> <value>

The commands are:

    ls                  list keys
    rm [key]            remove key "key"
    key                 get the value of "key"
    key value           overwrite the value of "key"

    clear               remove all values
    version             print application version"

`
