package main

import(
	"bytes"
	"fmt"
	"crypto/sha256"
)

type BlockChain struct{
	blocks []*Block
}

type Block struct{
	Hash [] byte
	Data [] byte
	PrevHash [] byte
}

func (b *Block) DeriveHash(){
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block{
	block:= &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string){
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

func Genesis() *Block{
	return CreateBlock("Team Kodo's genesis block", []byte{})
}

func InitBlockChain() *BlockChain{
	return &BlockChain{[]*Block{Genesis()}}
}

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