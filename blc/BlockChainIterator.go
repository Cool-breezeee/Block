package blc

import (
	"fmt"
	"github.com/boltdb/bolt"
)

// BlockChainIterator 迭代器
type BlockChainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

// Next 迭代器获得下一个区块
func (bci *BlockChainIterator) Next() *Block {
	var block *Block
	err := bci.DB.View(func(tx *bolt.Tx) error {
		// 1.表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockByte := b.Get(bci.CurrentHash)
			block = Deserialize(blockByte)
			bci.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return block
}
