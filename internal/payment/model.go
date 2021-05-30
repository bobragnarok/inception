package payment

type ReqPayment struct {
	Amount   int64  `json:"amount" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

type ReqInquiry struct {
	Status string `json:"status"`
}
