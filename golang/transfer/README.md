# Pegnet Transfer

A Pegnet transfer is the sending of a token from 1 address to another. Please read all the code and comments for how to do this in GoLang.

## Imports

The imports use the `Factom-Asset-Tokens` github organization, but pegnet is currently using a fork. Please do the following to ensure your imports are correct:

```bash
go mod edit -replace github.com/Factom-Asset-Tokens/factom=github.com/Emyrk/factom@rcd_full
go mod tidy
```