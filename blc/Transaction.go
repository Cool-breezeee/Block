package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

type Transaction struct {
	//  交易hash
	TxHash []byte
	//  输入
	VIns []*TXInput
	//  输出
	VOuts []*TXOutput
}

// NewCoinbaseTransaction 创世区块的Transaction
func NewCoinbaseTransaction(address string) *Transaction {
	// 消费
	txInput := &TXInput{
		TxHash:    []byte{},
		VOut:      -1,
		ScriptSig: "Genesis data",
	}
	// 未消费
	txOutput := &TXOutput{
		Value:        10,
		ScriptPubKey: address,
	}
	txCoinbase := &Transaction{
		TxHash: []byte{},
		VIns:   []*TXInput{txInput},
		VOuts:  []*TXOutput{txOutput},
	}
	// 设置Hash
	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// NewSimpleTransaction 交易的Transaction
func (bc *BlockChain) NewSimpleTransaction(from, to string, amount int64) *Transaction {
	//fmt.Println(unSpentTx)
	// 消费
	bytes, _ := hex.DecodeString("687b4303e08860424974964efae5562760967b206e1db9ad6f9834f8720c882f")

	txInput := &TXInput{
		TxHash:    bytes,
		VOut:      0,
		ScriptSig: from,
	}
	// 未消费
	txOutput1 := &TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	}
	txOutput2 := &TXOutput{
		Value:        10 - amount,
		ScriptPubKey: from,
	}
	txCoinbase := &Transaction{
		TxHash: []byte{},
		VIns:   []*TXInput{txInput},
		VOuts:  []*TXOutput{txOutput1, txOutput2},
	}
	// 设置Hash
	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) GetMoney() int64 {
	var money int64
	for _, out := range tx.VOuts {
		money += out.Value
	}
	return money
}
