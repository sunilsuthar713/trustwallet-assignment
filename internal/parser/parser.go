package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"trustwallet-assignment/pkg/models"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []models.Transaction
	PollNewBlocks()
}

type parser struct {
	subscribedAddresses map[string]bool
}

func NewParser() *parser {
	return &parser{
		subscribedAddresses: make(map[string]bool),
	}
}

func (p *parser) Subscribe(address string) bool {
	if _, exists := p.subscribedAddresses[address]; exists {
		return false // Already subscribed
	}
	p.subscribedAddresses[address] = true
	return true
}

func (p *parser) GetTransactions(address string) []models.Transaction {
	var transactions []models.Transaction

	// Iterate through the latest blocks for simplicity
	currentBlock := p.GetCurrentBlock()
	for i := currentBlock; i > currentBlock-10; i-- { // Check last 10 blocks
		blockHex := fmt.Sprintf("0x%x", i)

		// Prepare JSONRPC request for block details
		reqBody := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_getBlockByNumber",
			"params":  []interface{}{blockHex, true}, // true: Include full transaction objects
			"id":      1,
		}
		body, _ := json.Marshal(reqBody)

		resp, err := http.Post("https://ethereum-rpc.publicnode.com", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error fetching block:", err)
			continue
		}
		defer resp.Body.Close()

		var blockData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&blockData); err != nil {
			fmt.Println("Error decoding block data:", err)
			continue
		}

		// Extract transactions and filter by address
		if block, ok := blockData["result"].(map[string]interface{}); ok {
			if txns, ok := block["transactions"].([]interface{}); ok {
				for _, txn := range txns {
					tx := txn.(map[string]interface{})
					from, fromOk := tx["from"].(string)
					to, toOk := tx["to"].(string)
					hash, hashOk := tx["hash"].(string)

					// Skip transactions with missing fields
					if !fromOk || !toOk || !hashOk {
						fmt.Println("Skipping transaction due to missing fields")
						continue
					}

					transactions = append(transactions, models.Transaction{
						Hash:   hash,
						From:   from,
						To:     to,
						Amount: 0, // Amount logic can be added later
					})
				}
			}
		}
	}

	return transactions
}

func (p *parser) GetCurrentBlock() int {
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post("https://ethereum-rpc.publicnode.com", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error fetching current block:", err)
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected response status: %d\n", resp.StatusCode)
		return -1
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding JSONRPC response:", err)
		return -1
	}

	blockHex, ok := result["result"].(string)
	if !ok {
		fmt.Println("Invalid block number format in JSONRPC response")
		return -1
	}

	var blockNumber int
	fmt.Sscanf(blockHex, "0x%x", &blockNumber)
	return blockNumber
}

func (p *parser) PollNewBlocks() {
	lastParsedBlock := p.GetCurrentBlock()

	for {
		currentBlock := p.GetCurrentBlock()
		if currentBlock > lastParsedBlock {
			for i := lastParsedBlock + 1; i <= currentBlock; i++ {
				// Process each new block and filter transactions
				fmt.Printf("Processing block %d\n", i)
				for address := range p.subscribedAddresses {
					txns := p.GetTransactions(address)
					fmt.Printf("Transactions for %s: %+v\n", address, txns)
				}
			}
			lastParsedBlock = currentBlock
		}
		time.Sleep(5 * time.Second) // Poll every 5 seconds
	}
}
