package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/team-kodo/golang-blockchain.git/blockchain"
)

//CLI for the blockchain
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage() //no command entered so print instructions
		runtime.Goexit() //exits the app by shutting down the go routine, imp because badger needs to garbage collect K/V else can corrupt data
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block! ")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator() //convert blockchain to iter struct

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 { //reached genesis block as there is no PrevHash
			break
		}
	}
}

func (cli CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError) //CHECK
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] { //CHECK
	case "add":
		err := addBlockCmd.Parse(os.Args[2:]) //parse to check more than one command
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	//if there's no error
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage() //CHECK
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0) //fail safe to safely close the dB
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close() //defer only executes if the goExit is run properly

	cli := CommandLine{chain}
	cli.run()

}
