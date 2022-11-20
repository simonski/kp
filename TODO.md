# 1.0.0 TermUI

make an api
use the findfile
make some test harnesses

look at bubbletea

    Add TermUI interface.
    DONE -g mode or KP_GUI=1 to launch in gui termui mode
    DONE q quits
    a "good" ux for termui
    resize window
    include search functionality
    include demo functionality (prepopulate with lots of random keys to examine functionality)
    / for search
    asdw / jkli
    e edits
    enter selects into memory and quits
    fill to screen
    list keys and details, arrow keys, asdw and jkli for navigation
    d moves to delete/confirm
    mouse event over a thing


### Per-key passwords

```bash
kp put key -p
<password to store>
<password to encrypt with>
```

# 1.0.0

The road to 1.0 is finishing the TODO list.  Top of this page is the DOING section which
*should* be a single feature that would be 0.0.X incrementing on each feature.
Fix encryption

    - make own encryption key at startup via kp init
    - keep it simple (private, no public necessary)
    - kp init should do it
            
    - inspect certificate
        report on what it is, can we use it
    - initial setup - whenwe have no env vars etc introduce a ./kp setup
        look at different types of encryption keys as my regular one isn't working, weird
    - verify/docs
    - 'help' usage on each
    - 'update' command to update description but not value
        kp describe a "this is the thing"
    - fix typos, docs, help, 
    - retain history as older copies

    - move it all out to a library so I can use it externally
        github.com/simonski/kp
            functions as KP_xxxx

    - DONE a changelog
    - DONE move to a crypto package or utils for other usage
    - DONE move cli to a cli project from goutils?

## DOING

branch: features/prep_v1

- DONE: basics of encryption work.
- Improve error handling and information.
- "init", "verify", "setup" steps to be described and completed
- "info" is almost pointless
- `kp verify` ssh public key verification/generation

- investigate overall binary size (5MB - can it be smaller?)

## TODO

Before I can release on @master and @1.0.0 I need to

- review the README/create a USER_GUIDE ./kp help
- attempt to format the STDOUT based on $COLUMNS or tput if possible?
- review unit tests/coverage
- setup dependency on goutils as formal - version goutils as github.com/simonski/goutils@1.0.0
- verify initial setup - mac
- verify initial setup - win
- verify initial setup - linux
- "help" on all function calls
- github actions from crosscompiling
- verify minimum go compatible version
- ?move to sqlite?

## 0.0.2  IN PROGRESS

- prettyprint: use large descriptions etc, so when printing to stdout

## 0.0.1

- ~~"I don't know how to 'x' did you mean 'put'?"~~
- ~~get/put stdout~~
- ~~accept input secretly by default via  crpytic put~~
- ~~description on key~~
- ~~created, last_updated~~
- ~~"description" in keys~~
- ~~better stdout formatting (look at padding)~~
- ~~decide on set/get/put and reading from stdin~~
