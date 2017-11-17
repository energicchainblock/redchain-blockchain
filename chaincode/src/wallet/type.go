package wallet

type Wallet struct {
	Available float64 `json:"available"`
	Ico       float64
	Version   int64
}

type CommonReq struct {
	Number float64 `json:"number"`
}

// type CommonReply struct {
// 	Address   string  `json:"address"`
// 	Available float64 `json:"available"`
// }

type CommonReply struct {
	Address string  `json:"address"`
	Value   float64 `json:"value"`
	From    string  `json:"from"`
	Fvalue  float64 `json:"fvalue"`
	Amount  float64 `json:"amount"`
}
