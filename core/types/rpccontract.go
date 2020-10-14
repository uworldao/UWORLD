package types

type RpcContract struct {
	Contract    string               `json:"contract"`
	Name        string               `json:"name"`
	Abbr        string               `json:"abbr"`
	Increase    bool                 `json:"increase"`
	Description string               `json:"description"`
	Records     []*RPcContractRecord `json:"records"`
}

type RPcContractRecord struct {
	Height   uint64  `json:"height"`
	TxHash   string  `json:"txhash"`
	Time     uint64  `json:"time"`
	Amount   float64 `json:"amount"`
	Receiver string  `json:"receiver"`
}

func TranslateContractToRpcContract(contract *Contract) *RpcContract {
	rpcContract := &RpcContract{
		Contract:    contract.Contract,
		Name:        contract.CoinName,
		Abbr:        contract.CoinAbbr,
		Increase:    contract.IncreaseSwitch,
		Description: contract.Description,
		Records:     make([]*RPcContractRecord, contract.Records.Len()),
	}
	for i, record := range *contract.Records {
		rpcContract.Records[i] = &RPcContractRecord{
			Height:   record.Height,
			TxHash:   record.TxHash.String(),
			Time:     record.Time,
			Amount:   Amount(record.Amount).ToCoin(),
			Receiver: record.Receiver,
		}
	}
	return rpcContract
}
