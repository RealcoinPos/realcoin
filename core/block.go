package core

import (
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index		int		`json:"index"`
	Timestamp	int		`json:"timestamp"`	
	Data		string	`json:"data"`
	Hash		string	`json:"hash"`
	PrevHash	string	`json:"prevhash"`
	Nonce		int		`json:"nonce"`
}


var posnum=0
var Blockchain []Block

func calculateHash(block Block) string {
	record := string(block.Index) + string(block.Timestamp) + block.Data + block.PrevHash + string(block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenesisChain() []Block {
	gBlock := genesisBlock("Genesis Block Start")
    gChain := append(Blockchain, gBlock)
	return gChain
}


func genesisBlock(genesisData string) Block {
	var block=*new(Block)
	bFind:=blockFinder("0_")
	if bFind {
		return block
	}
	block.Index=0
	block.Timestamp=int(time.Now().Unix())
	block.Nonce=GetNonce()
	block.Data=genesisData
	block.PrevHash="0"
	block.Hash=calculateHash(block)
	blockFile(block)
	return block
}

func generateBlock(addData string) Block{
	var block Block
	var oblock Block
	c:=CurrentHeight()
	if c<1 {
		return block
	}
	bname:=blockName(strconv.Itoa(c-1))
	err:=fileRead(bname, oblock)
	Err(err,0)
	block.Index=c
	block.Timestamp=int(time.Now().Unix())
	block.Nonce=GetNonce()
	block.Data=addData
	block.PrevHash=oblock.Hash
	block.Hash=calculateHash(block)
	return block
}

func addBlock(data string) {
	var oBlock=Blockchain[len(Blockchain)-1];
	var nBlock=generateBlock(data);
	if isBlockValid(nBlock, oBlock) {
		blockFile(nBlock)
		newBlockchain := append(Blockchain, nBlock)
		replaceChain(newBlockchain)
	}
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {	return false }
	if oldBlock.Hash != newBlock.PrevHash {	return false }
	if calculateHash(newBlock) != newBlock.Hash { return false }
	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func GetNonce() int {
	return posnum
}