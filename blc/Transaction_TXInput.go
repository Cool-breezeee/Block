package blc

type TXInput struct {
	// 交易 hash
	TxHash []byte
	// 存储TXOutput在VOut里面的索引
	VOut int
	/// 用户名
	ScriptSig string
}

// UnLockWithAddress 判断当前的消费是否是用户的钱
func (txInput *TXInput) UnLockWithAddress(address string) bool {
	return txInput.ScriptSig == address
}

func (tx *Transaction) IsCoinBaseTransaction() bool {
	return len(tx.VIns[0].TxHash) == 0 && tx.VIns[0].VOut == -1
}
