package main

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{createGenesisBlock()},
	}
}

func createGenesisBlock() *Block {
	return NewBlock(0, "Genesis Block", "")
}

func (bc *Blockchain) AddBlock(data string) {
	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(previousBlock.Index-1, data, previousBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}