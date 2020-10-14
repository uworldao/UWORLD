package types

type RpcContractTransactionBody struct {
	Contract    string `json:"contract"`
	To          string `json:"to"`
	Name        string `json:"name"`
	Abbr        string `json:"abbr"`
	Description string `json:"description"`
	Increase    bool   `json:"increase"`
	Amount      uint64 `json:"amount"`
}
