package blc

import (
	"Block/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 256位Hash里面前面至少要有16个0
const targetBit = 16

type ProofOfWork struct {
	Block  *Block   //需要验证的区块
	target *big.Int //大数据存储
}

// NewProofOfWork 创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	// 1. 创建初始值为1的target
	target := big.NewInt(1)
	// 2. 左移256 - targetBit 位
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.HashTransactions(),
			pow.Block.PrevBlockHash,
			utils.Int64ToBytes(pow.Block.BlockHeight),
			utils.Int64ToBytes(pow.Block.Timestamp),
			utils.Int64ToBytes(nonce),
			utils.Int64ToBytes(int64(targetBit)),
		}, []byte{},
	)
	return data
}

func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.Block.Hash)
	// 1. pow.Block.Hash <-> pow.target
	if pow.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

func (pow *ProofOfWork) Run() ([]byte, int64) {
	nonce := 0
	var hashInt big.Int // 存储新生成的hash值
	var hash [32]byte
	for {
		// 1. 将Block属性拼接成字节数组
		dataByte := pow.prepareData(int64(nonce))
		// 2. 生成Hash
		hash = sha256.Sum256(dataByte)
		fmt.Printf("\r%x", hash)
		// 3. 将hash存入hashInt
		hashInt.SetBytes(hash[:])
		// 4. 判断hashInt是否小于Block的target，满足条件，跳出循环
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		// 生成的比目标小，说明转换成二进制开头的0的数量符合要求
		if pow.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce += 1
	}
	fmt.Println()
	return hash[:], int64(nonce)
}
