package blc

type TXOutput struct {
	Value        int64
	ScriptPubKey string //用户名
}

// UnLockWithAddress 判断当前的消费是谁的
func (txOutput *TXOutput) UnLockWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
