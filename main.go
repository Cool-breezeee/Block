package main

import "Block/cli"

const (
	dbName         = "blockchain.db"
	blockTableName = "test"
)

func main() {
	cli := cli.CLI{}
	cli.Run()
}
