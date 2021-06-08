package main

import (
	"fmt"

	"github.com/team-kodo/golang-blockchain.git/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("First")
	chain.AddBlock("Second")
	chain.AddBlock("Third")

	for _, block := range chain.Blocks {

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

	}
}