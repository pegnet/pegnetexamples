package main

// !!!		PLEASE READ		!!!
// ---		Import Advice 		---
// Pegnet has plans to integrate into the FAT ecosystem, but not at this time.
// Because Pegnet is currently it's own project and ecosystem, it has outpaced
// the FAT dependencies in some areas.
// If you are using GoLang to create transactions, we recommend using the
// `go mod` native dependency manager. Your imports should look like the imports
// in this file, but you will need to edit your go.mod file to replace
// the FAT path with the fork pegnet is using. This can be done with the following
// go mod command:
//
// go mod edit -replace github.com/Factom-Asset-Tokens/factom=github.com/Emyrk/factom@rcd_full
// go mod tidy
//
// This will use the right dependencies, but also keep your imports correct.
// github.com/Factom-Asset-Tokens/factom is the upstream repo all changes in
// the fork will eventually roll into.

import (
	"fmt"

	"github.com/Factom-Asset-Tokens/factom"
	"github.com/pegnet/pegnetd/fat/fat2"
)

const (
	// Address 'FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q' has
	// FCT on a local/custom dev net
	FCTPrivateAddress = "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"

	// Recipient is who the transfer is going to.
	// The private address for this address is
	// 'Fs39xgC7ZGyTnTUnM2zgjfn28PvM2QA2DHhkHNBUqLMVFKpdqtqQ'
	Recipient = "FA3ppia1ZbcuXHdhAYLKBU89JKNVdFJVGiRQS3ziNtqhcoFvFxqJ"

	// Public address is "EC3TsJHUs8bzbbVnratBafub6toRYdgzgbR7kWwCW4tqbmyySRmg"
	ECPayment = "Es2XT3jSxi1xqrDvS5JERM3W3jh1awRHuyoahn3hbQLyfEi1jvbq"
)

var (
	// This is the pegnet tx chain. It is a constant that we need to include
	// to submit our entry to the right chain in Factom.
	TransactionChain = factom.NewBytes32("cffce0f409ebba4ed236d49d89c70e4bd1f1367d86402a3363366683265a242d")
)

func main() {
	// If your factom is not at localhost, you will need to adjust this client
	cl := factom.NewClient()
	cl.FactomdServer = "http://localhost:8088/v2"

	// This code will provide an example for how to craft a pegnet transfer.
	// The transfer will be of 100 PEG to the 'Recipient'.
	sender, _ := factom.NewFsAddress(FCTPrivateAddress)
	recipient, _ := factom.NewFAAddress(Recipient)
	payment, _ := factom.NewEsAddress(ECPayment)

	// All pegnet transactions can contain 1 or more transfers/conversions.
	// For the purposes of this demo, we will only create 1 transfer, but
	// you could easily combine multiple transfers in the same entry.
	var batch fat2.TransactionBatch
	batch.Version = 1
	batch.Entry.ChainID = &TransactionChain

	// Let's craft the actual transfer
	var transfer fat2.Transaction
	// Input
	// 1 PEG is the amount 1e8
	transfer.Input.Amount = 100 * 1e8           // 100 PEG
	transfer.Input.Address = sender.FAAddress() // Sender address
	transfer.Input.Type = fat2.PTickerPEG       // PEG Currency

	// Output (Only 1)
	transfer.Transfers = make([]fat2.AddressAmountTuple, 1)
	transfer.Transfers[0].Amount = 100 * 1e8 // 100 PEG (matches inputs)
	transfer.Transfers[0].Address = recipient

	// Add to the batch
	batch.Transactions = []fat2.Transaction{transfer}

	// Transaction is crafted, time to sign and get the Factom Entry we
	// can submit
	entry, err := batch.Sign(sender)
	if err != nil {
		panic(err)
	}

	balance, _ := payment.GetBalance(nil, cl)
	if balance == 0 {
		// You need entry credits to send a pegnet transaction.
		// If you do not have any, the entry will not be submitted.
		panic("no entry credits")
	}

	commit, err := entry.ComposeCreate(nil, cl, payment)
	if err != nil {
		panic(err)
	}

	fmt.Println("Transaction Sent")
	// You can check the 'Entry Hash' on the factom explorer
	fmt.Printf("Entry Hash : %s\n", entry.Hash.String())
	// You can check the txid on the pExplorer
	fmt.Printf("Txid       : 0-%s\n", entry.Hash.String())
	fmt.Printf("Commit Hash: %s\n", commit.String())
}
