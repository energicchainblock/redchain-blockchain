package wallet

type Wallet struct {
	Balance float64 `json:"balance"`
	In      float64 `json:"in"`
	Out     float64 `json:"out"`
}

type RechargeReq struct {
	Way   int8    `json:"way"`
	Value float64 `json:"value"`
}

type TransferReq struct {
	Value float64 `json:"value"`
}
