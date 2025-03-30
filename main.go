package main

import (
	"fmt"
	"trustwallet-assignment/internal/parser" // Import the parser package
)

func main() {
	// Create a new parser instance
	p := parser.NewParser()

	// Fetch the current block and display it
	currentBlock := p.GetCurrentBlock()
	fmt.Println("Current Block:", currentBlock)

	// Subscribe to a test Ethereum address
	address := "0xAbc123..." // Replace with a valid Ethereum address
	subscribed := p.Subscribe(address)
	if subscribed {
		fmt.Printf("Subscribed to address: %s\n", address)
	} else {
		fmt.Printf("Address already subscribed: %s\n", address)
	}

	// Fetch and display transactions for the subscribed address
	transactions := p.GetTransactions(address)
	fmt.Printf("Transactions for %s:\n", address)
	for _, txn := range transactions {
		fmt.Printf("- Hash: %s, From: %s, To: %s, Amount: %d\n", txn.Hash, txn.From, txn.To, txn.Amount)
	}

	// Start polling for new blocks dynamically
	fmt.Println("Starting to poll for new blocks...")
	p.PollNewBlocks()
}
