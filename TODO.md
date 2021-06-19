# 1.0.0

The road to 1.0 is finishing the TODO list.  Top of this page is the DOING section which
*should* be a single feature that would be 0.0.X incrementing on each feature.

## DOING

- "init", "verify", "setup" steps to be described and completed
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
