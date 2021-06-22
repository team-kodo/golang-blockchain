package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Team Kodo's genesis block", []byte{})
}

/**
function to convert block into bytes as BadgerDB can only store bytes
**/
func (b *Block) Serialize() []byte {
	var res bytes.Buffer            //Store bytes in buffer memory
	encoder := gob.NewEncoder(&res) //create new encoder

	err := encoder.Encode(b) //returns an err

	Handle(err)

	return res.Bytes() //returns byte output of the block
}

func Deserialize(data []byte) *Block {
	var block Block //new block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block //return decoded block

}

//error handling function
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
