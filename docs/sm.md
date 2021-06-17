% sm 8
# NAME
sm - AWS Secrets Manager CLI Tool
# SYNOPSIS
sm


# COMMAND TREE

- [get, view](#get-view)
- [edit, e](#edit-e)
- [create, c](#create-c)
- [put](#put)
- [delete, del](#delete-del)
- [list](#list)
- [describe](#describe)

**Usage**:
```
sm [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## get, view

select from list or pass in specific secret

**--binary, -b**: get the SecretBinary value

**--secret-id, -s**="": Specific Secret to view, will bypass select/search

## edit, e

interactive edit of a secret String Value

**--binary, -b**: get the SecretBinary value

**--secret-id, -s**="": Specific Secret to edit, will bypass select/search

## create, c

create new secret in Secrets Manager

**--binary, -b**: get the SecretBinary value

**--description, -d**="": Additional description text.

**--interactive, -i**: Open interactive editor to create secret value. If no 'value' is provided, an editor will be opened by default.

**--secret-id, -s**="": Secret name

**--tags, -t**="": key=value tags (CSV list)

**--value, -v**="": Secret Value. Will store as a string, unless binary flag is set.

## put

non-interactive update to a specific secret

```
Stores a new encrypted secret value in the specified secret. To do this, the 
operation creates a new version and attaches it to the secret. The version 
can contain a new SecretString value or a new SecretBinary value.

This will put the value to AWSCURRENT and retain one previous version 
with AWSPREVIOUS.
```

**--binary, -b**: get the SecretBinary value

**--interactive, -i**: Override and open interactive editor to verify and modify the new secret value.

**--secret-id, -s**="": Secret name

**--value, -v**="": Secret Value. Will store as a string, unless binary flag is set.

## delete, del

delete a specific secret

**--force, -f**: Bypass recovery window (30 days) and immediately delete Secret.

**--secret-id, -s**="": Specific Secret to delete

## list

display table of all secrets with meta data

## describe

print description of secret to `STDOUT`

**--secret-id, -s**="": Specific Secret to describe, will bypass select/search

