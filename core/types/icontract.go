package types

type IContract interface {
	AddContract(height uint64, txHash string, time uint64, amount uint64) bool
	IsExist(txId string) bool
	FallBack(height uint64) error
}
