package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"time"
)

type Block struct {
	Hash      []byte
	Data      []byte
	PrevHash  []byte
	Timestamp time.Time
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash, []byte(b.Timestamp.String())}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func (b *Block) IsValid() bool {
	info := bytes.Join([][]byte{b.Data, b.PrevHash, []byte(b.Timestamp.String())}, []byte{})
	hash := sha256.Sum256(info)
	return bytes.Equal(hash[:], b.Hash)
}

type BlockChain struct {
	blocks []*Block
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, time.Now()}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func (chain *BlockChain) IsValidChain() bool {
	for i := 1; i < len(chain.blocks); i++ {
		if !chain.blocks[i].IsValid() || !bytes.Equal(chain.blocks[i-1].Hash, chain.blocks[i].PrevHash) {
			return false
		}
	}
	return true
}

func (chain *BlockChain) PrintBlocks() {
	for _, block := range chain.blocks {
		fmt.Printf("Timestamp: %s\n", block.Timestamp.Format(time.RFC3339))
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("--------------------------------------")
	}
}

func main() {
	chain := InitBlockChain()

	// Adicione blocos com dados fornecidos pela linha de comando
	for i, arg := range os.Args[1:] {
		chain.AddBlock(fmt.Sprintf("Block %d: %s", i+1, arg))
	}

	if chain.IsValidChain() {
		chain.PrintBlocks()
	} else {
		fmt.Println("Blockchain is not valid.")
	}
}
