package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

//struct to access each block
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

//updated to store blockchain in the DB
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	//connect db
	db, err := badger.Open(opts) //outputs pointer to db and err
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error { //txn is transaction
		//Check blockchain exists
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")

			genesis := Genesis()
			fmt.Println("Genesis proved")

			err = txn.Set(genesis.Hash, genesis.Serialize()) //Set genesis as serialized bytes in the DB
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash //storing lasthash in memory
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value()
			return err
		}

	})
	Handle(err)
	blockchain := BlockChain{lastHash, db} //add blockchain to memory
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	//read only type transaction
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh")) //get last hash from DB
		Handle(err)
		lastHash, err = item.Value()

		return err
	})

	Handle(err)
	newBlock := CreateBlock(data, lastHash) //create new block with data passed to func and with lasthash

	//add new block to db and assign new block's hash as last hash
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash //get blockchain and set the hash of the new block as last hash

		return err
	})
}

//blockchain iterator to go over each block
func (chain *BlockChain) iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

//go through each block; from end to genesis as we have the last hash value
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error { //badger's iter func
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
