package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const passPhrase = "blockchain"

type Block struct {
	Data       []byte
	Nonce      int
	Difficulty int
	Timestamp  int64
	Hash       []byte
	PrevHash   []byte
	//TransactionData []byte
}

type Blockchain struct {
	blocks []*Block
}

type Transaction struct {
	Sender   string
	Receiver string
	Amount   int
	// Timestamp int64
}

//type Transactions struct {
//	AllTransactions []*Transaction
//}

func (t *Transaction) JsonToByte() []byte {
	trans, err := json.Marshal(t)
	if err != nil {
		log.Print(err)
	}
	return trans
}

func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (b *Block) EncryptData(data []byte, passPhrase string) {
	block, _ := aes.NewCipher([]byte(CreateHash(passPhrase)))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	b.Hash = cipherText[:]
}

func CreateBlock(ts []byte, prevHash []byte, passPhrase string) *Block {
	t := time.Now().UnixNano()
	block := &Block{Data: ts, PrevHash: prevHash, Timestamp: t, Nonce: 1, Difficulty: 2, Hash: []byte{}}
	block.EncryptData(bytes.Join([][]byte{block.Data, block.PrevHash}, []byte{}), passPhrase)
	return block
}

func (chain *Blockchain) AddBlock(ts *Transaction, passPhrase string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(ts.JsonToByte(), prevBlock.Hash, passPhrase)
	chain.blocks = append(chain.blocks, newBlock)
}

func Decrypt(data []byte, passPhrase string) []byte {
	key := []byte(CreateHash(passPhrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Print(err)
		panic(err.Error())

	}
	return plaintext
}

func (b *Block) ByteToJson(data []byte) any {
	var newData []interface{}
	_ = json.Unmarshal(data, &newData)
	return newData
}

func GenesisBlock() *Block {
	ts := &Transaction{Sender: "", Receiver: "", Amount: 0}
	return CreateBlock(ts.JsonToByte(), []byte{}, passPhrase)
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func GenerateGenesisBlock(w http.ResponseWriter, r *http.Request) {
	gen := InitBlockchain()
	for _, block := range gen.blocks[:1] {
		fmt.Fprintf(w, "Hash : %x\n", block.Hash)
		fmt.Fprintf(w, "Previous Hash : %x\n", block.PrevHash)
		data := Decrypt(block.Hash, passPhrase)
		block.ByteToJson(data)
		w.Write(data)
		// fmt.Fprintf(w, "Data : %s\n", data)
	}
	// fmt.Fprintf(w, "Genesis Block: %v\n", gen)
}

func MainChain(w http.ResponseWriter, r *http.Request) {
	chain := InitBlockchain()
	ts := &Transaction{Sender: "raihan", Receiver: "shuvo", Amount: 100}
	chain.AddBlock(ts, passPhrase)
	for _, block := range chain.blocks {
		fmt.Fprintf(w, "Hash : %x\n", block.Hash)
		fmt.Fprintf(w, "Previous Hash : %x\n", block.PrevHash)
		data := Decrypt(block.Hash, passPhrase)
		block.ByteToJson(data)
		fmt.Fprintf(w, "Data : %s\n", data)
		fmt.Fprintf(w, "%s\n", strings.Repeat("=", 50))

	}

}
