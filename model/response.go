package model

type PrintChainResponse struct {
	Height       int64
	Nonce        int64
	PreBlockHash string
	Hash         string
	Ts           string
	Tx           []Transaction
}
type Transaction struct {
	TxHash string
}
