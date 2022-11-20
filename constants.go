package main

import "strings"

// KP_FILE the key for the env var pointint to the file we load/save
const KP_FILE = "KP_FILE"

// KP_KEY the encypt/decrypt key
const KP_KEY = "KP_KEY"

const DEFAULT_KEY_FILE = "~/.ssh/kp.id_rsa"
const DEFAULT_DB_FILE = "~/.kpfile"

// GLOBAL_USAGE - well, it tells me what to type
const GLOBAL_USAGE = `kp is a tool for using key/pairs.

Usage:

    kp <command> [arguments]

The commands are:

    ls  <search term>               list keys 
         -a (include hidden)

    put <key>                       save "key/value" (read stdin for value)

    get <key> (-stdout)             retrieve key/value to clipboard (or -stdout)

    update <key>                    update metadata on the key
         -description                   
         -type
         -url
         -username
         -note

    open <key>                      opens the url associated with the key 

    rename <key1> <key2>            rename "key1" to "key2"
    rm <key>                        permanently remove "key"

    encrypt <value>                 encrypt the value 
    decrypt <value>                 decrypt the value

    tag   <key> <tag>               add tag to a key
    untag <key> <tag>               remove tag from a key

    hide <key>                      archive (hide) the key
    show <key>                      unarchive (make visible) the key

    info                            review environment variables used
    verify                          check encryption keys exist and work
    version                         print application version

`

const GLOBAL_SSH_KEYGEN_USAGE = `The following will create a suitable encryption key: 

     TOKEN_DEFAULT_SSH_COMMAND

You can optionally use environment variables to override the defaults:

     export KP_FILE=TOKEN_DEFAULT_DB_FILE
     export KP_KEY=TOKEN_DEFAULT_KEY_FILE

`

const DEFAULT_OPENSSH_COMMAND = "ssh-keygen -b 2048 -t rsa -N \"\" -m pkcs8 -f TOKEN_DEFAULT_KEY_FILE"

func GetSSHCommand(key string) string {
	return strings.ReplaceAll(DEFAULT_OPENSSH_COMMAND, "TOKEN_DEFAULT_KEY_FILE", key)
}
