package blc

import (
	"Block/model"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

const (
	dbName         = "blockchain.db"
	blockTableName = "test"
)

type BlockChain struct {
	Tip []byte //最新的区块的Hash
	DB  *bolt.DB
}

// DBExists 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateNewBlockChainWithGenesisBlock 创建带有创世区块的区块链
func CreateNewBlockChainWithGenesisBlock(address string) *BlockChain {
	if DBExists() {
		fmt.Println("创世区块已经存在...")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块...")
	// 创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var genesisHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			if err.Error() != "bucket already exists" {
				return err
			}
		}
		if b != nil {
			//创建一个Coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address)
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到表中
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最近的区块的Hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return &BlockChain{
		Tip: genesisHash,
		DB:  db,
	}
}

// Iterator 通过 BlockChain 生成迭代器对象
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		CurrentHash: bc.Tip,
		DB:          bc.DB,
	}
}

// AddNewBlockToBlockChain 添加新的区块到区块链中
func (bc *BlockChain) AddNewBlockToBlockChain(txs []*Transaction) error {

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 1 getTable
		b := tx.Bucket([]byte(blockTableName))
		// 2 create new block
		if b != nil {
			// use Tip to get the newest Block
			blockBytes := b.Get(bc.Tip)
			block := Deserialize(blockBytes)
			newBlock := CreateNewBlock(txs, block.BlockHeight+1, block.Hash)
			// 3 serialize block and save in database
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Println(err)
				return err
			}
			// 4 update table l hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Println(err)
				return err
			}
			// 5 update blockchain tip
			bc.Tip = newBlock.Hash
		} else {
			return errors.New("db is err")
		}
		return nil
	})
	return err
}

// PrintChain 遍历输出所有区块的信息
func (bc *BlockChain) PrintChain() {
	blockChainIterator := bc.Iterator()
	for {
		block := blockChainIterator.Next()
		var prevBlock model.PrintChainResponse
		prevBlock.Height = block.BlockHeight
		prevBlock.PreBlockHash = fmt.Sprintf("%x\n", block.PrevBlockHash)
		prevBlock.Hash = fmt.Sprintf("%x\n", block.Hash)
		prevBlock.Ts = time.UnixMilli(block.Timestamp).Format("2006-01-02 15:04:05.000")
		prevBlock.Nonce = block.Nonce
		//fmt.Printf("Height = %d\n", block.BlockHeight)
		//fmt.Printf("PreBlockHash = %x\n", block.PrevBlockHash)
		//fmt.Printf("Ts = %s\n", time.UnixMilli(block.Timestamp).Format("2006-01-02 15:04:05.000"))
		//fmt.Printf("Hash = %x\n", block.Hash)
		//fmt.Printf("Nonce = %d\n", block.Nonce)
		//fmt.Println("Txs:")
		for _, tx := range block.Txs {
			var prevTx model.Transaction
			prevTx.TxHash = fmt.Sprintf("%x\n", tx.TxHash)
			//fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("VIns = ")
			for _, in := range tx.VIns {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.VOut)
				fmt.Printf("%s\n", in.ScriptSig)
			}
			fmt.Println("VOuts = ")
			for _, out := range tx.VOuts {
				fmt.Printf("%d\n", out.Value)
				fmt.Printf("%s\n", out.ScriptPubKey)
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
		fmt.Println("------------------")
	}
}

func GetBlockChainObject() *BlockChain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	return &BlockChain{
		Tip: tip,
		DB:  db,
	}
}

// UnUTXos 如果一个地址所对应的TXOuuPut未花费，那么这个Transaction应该添加到数组返回
func (bc *BlockChain) UnUTXos(address string) []*TXOutput {
	var (
		unSpentTxs []*TXOutput
	)
	spentTxOutputs := make(map[string][]int)

	blockIterator := bc.Iterator()
	for {
		block := blockIterator.Next()
		for _, tx := range block.Txs {
			if !tx.IsCoinBaseTransaction() {
				for _, in := range tx.VIns {
					if in.UnLockWithAddress(address) {
						hash := hex.EncodeToString(in.TxHash)
						spentTxOutputs[hash] = append(spentTxOutputs[hash], in.VOut)
					}
				}
			}
			for index, out := range tx.VOuts {
				if out.UnLockWithAddress(address) {
					if spentTxOutputs != nil {
						if indexArray, ok := spentTxOutputs[hex.EncodeToString(tx.TxHash)]; ok {
							for _, value := range indexArray {
								if value == index {
									goto Loop
								}
							}
							unSpentTxs = append(unSpentTxs, out)
						} else {
							unSpentTxs = append(unSpentTxs, out)
						}
					}
				}
			Loop:
				continue
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return unSpentTxs
}

// MineNewBlock 挖掘新的区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	// ./main send -from '[\"jayson.sun\"]' -to '[\"tao.liu2\"]' -amount '[\"2\"]'
	// [jayson.sun]
	// [tao.liu2]
	// [2]

	// 1.通过相关算法建立交易 Transaction 数组
	var txs []*Transaction
	// 建立交易
	pay, _ := strconv.Atoi(amount[0])
	tx := bc.NewSimpleTransaction(from[0], to[0], int64(pay))
	txs = append(txs, tx)

	var block *Block
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = Deserialize(blockBytes)
		}
		return nil
	})
	// 2.建立新的区块
	block = CreateNewBlock(txs, block.BlockHeight+1, block.Hash)

	// 3.更新存储数据库
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("l"), block.Hash)
			bc.Tip = block.Hash
		}
		return nil
	})
}
