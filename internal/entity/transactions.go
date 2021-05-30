package entity

type Transactions struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Token    string `json:"token"`
	Status   string `json:"status"`
}
