package main

import (
	"fmt"
)

func main(){
	chain:= InitBlockChain()

	chain.AddBlock("Poopoo")
	chain.AddBlock("pEEPEE")
	chain.AddBlock("WEEWEE")

	for _, block := range chain.blocks{
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Block's hash: %x\n\n", block.Hash)
	}
}