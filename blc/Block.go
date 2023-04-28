package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// Block 区块
type Block struct {
	//1. 区块高度
	BlockHeight int64
	//2. 上一个区块的Hash
	PrevBlockHash []byte
	//3. 交易数据
	Txs []*Transaction
	//4. 时间戳
	Timestamp int64
	//5. 当前区块的Hash
	Hash []byte
	//6. 工作量Nonce
	Nonce int64
}

// Serialize 将区块序列化成字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// Deserialize 将字节数组反序列化成区块
func Deserialize(byte []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(byte))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// CreateNewBlock 创建新的区块
func CreateNewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	block := &Block{
		BlockHeight:   height,
		PrevBlockHash: prevBlockHash,
		Txs:           txs,
		Timestamp:     time.Now().UnixMilli(),
		Hash:          nil,
		Nonce:         0,
	}
	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)
	// 挖矿验证
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// SetBlockHash 获得区块的Hash值
//func (b *Block) SetBlockHash() {
//	// 区块高度转字节数组
//	heightBytes := utils.Int64ToBytes(b.BlockHeight)
//	// 时间戳转字节数组
//	timeBytes := utils.Int64ToBytes(b.Timestamp)
//	// 拼接字节数组
//	blockBytes := bytes.Join([][]byte{heightBytes, b.PrevBlockHash, b.Txs, timeBytes, b.Hash}, []byte{})
//	// 求当前区块的hash值
//	hash := sha256.Sum256(blockBytes)
//	// 将hash赋值给当前区块
//	b.Hash = hash[:]
//}

// CreateGenesisBlock 创建创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return CreateNewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

// HashTransactions 将Txs转换成字节数组
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
