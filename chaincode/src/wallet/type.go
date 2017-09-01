package wallet

type Wallet struct {
	Available float64 `json:"available"`
	Ico       float64 `json:"ico"`
	Version   int64   `json:"version"`
}

type CommonReq struct {
	Number float64 `json:"number"`
}

type CommonReply struct {
	Address   string  `json:"address"`
	Available float64 `json:"available"`
}
