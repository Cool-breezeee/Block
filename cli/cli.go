package cli

import (
	"Block/blc"
	"Block/utils"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

type cli struct {
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tsend -from FROM -to To -amount Amount 交易明细")
	fmt.Println("\tcreateBlockChain -address 添加DATA为创世区块的区块链")
	fmt.Println("\tprintChain -- 输出区块信息")
	fmt.Println("\tgetBalance -address 查询账户余额")
}

type CLI interface {
	Transaction() gin.HandlerFunc
	PrintChain() gin.HandlerFunc
	GetBalance() gin.HandlerFunc
	CreateBlockChain() gin.HandlerFunc
}

func (cli *cli) PrintChain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		nodeId := ctx.Param("node_id")
		fmt.Println(nodeId)
		cli.printChain()
	}
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *cli) addBlock(txs []*blc.Transaction) {
	if blc.DBExists() {
		blockChain := blc.GetBlockChainObject()
		defer blockChain.DB.Close()
		blockChain.AddNewBlockToBlockChain(txs)
	} else {
		fmt.Println("请先调用 createBlockChain 创建区块链数据库")
		os.Exit(1)
	}
}
func (cli *cli) createBlockChain(address string) {
	blockChain := blc.CreateNewBlockChainWithGenesisBlock(address)
	defer blockChain.DB.Close()
}

func (cli *cli) getBalance(address string) {
	fmt.Println("地址 = ", address)
	blockChain := blc.GetBlockChainObject()
	defer blockChain.DB.Close()
	txOutputs := blockChain.UnUTXos(address)
	for _, out := range txOutputs {
		fmt.Println(out)
	}
}

// 转账
func (cli *cli) send(from, to, amount []string) {
	if !blc.DBExists() {
		fmt.Println("请先调用 createBlockChain 创建区块链数据库")
		os.Exit(1)
	}
	blockChain := blc.GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.MineNewBlock(from, to, amount)
}

func (cli *cli) printChain() {
	if blc.DBExists() {
		blockChain := blc.GetBlockChainObject()
		defer blockChain.DB.Close()
		blockChain.PrintChain()
	} else {
		fmt.Println("请先调用 createBlockChain 创建区块链数据库")
		os.Exit(1)
	}

}

func (cli *cli) Run() {
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	getBalanceChainCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	flagSendFromData := sendCmd.String("from", "", "转账源地址")
	flagSendToData := sendCmd.String("to", "", "转账目的地址")
	flagSendAmountData := sendCmd.String("amount", "", "转账金额")
	createBlockChainData := createBlockChainCmd.String("address", "", "创建创世区块的地址")
	getBalanceData := getBalanceChainCmd.String("address", "", "查询余额")
	isValidArgs()
	switch os.Args[1] {
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
	case "createBlockChain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
	case "getBalance":
		err := getBalanceChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		printUsage()
		os.Exit(1)
	}
	if sendCmd.Parsed() {
		if *flagSendFromData == "" || *flagSendToData == "" || *flagSendAmountData == "" {
			printUsage()
			os.Exit(1)
		} else {
			from := utils.JsonToArray(*flagSendFromData)
			to := utils.JsonToArray(*flagSendToData)
			amount := utils.JsonToArray(*flagSendAmountData)
			cli.send(from, to, amount)
		}
	}
	if createBlockChainCmd.Parsed() {
		if *createBlockChainData == "" {
			fmt.Println("地址不能为空!")
			printUsage()
			os.Exit(1)
		} else {
			cli.createBlockChain(*createBlockChainData)
		}
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if getBalanceChainCmd.Parsed() {
		if *getBalanceData == "" {
			fmt.Println("地址不能为空!")
			printUsage()
			os.Exit(1)
		} else {
			cli.getBalance(*getBalanceData)
		}
	}
}
