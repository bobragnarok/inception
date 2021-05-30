package payment

import (
	"inception/internal/entity"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/spf13/viper"
)

type Service interface {
	Payment(req ReqPayment) error
	Inquiry(status string) ([]entity.Transactions, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return service{repository}
}

func (s service) Payment(req ReqPayment) error {
	client, err := omise.NewClient(viper.GetString("omise.public"), viper.GetString("omise.key"))
	if err != nil {
		return err
	}

	token, retrieve := &omise.Token{}, &operations.RetrieveToken{
		ID: "tokn_test_4xs9408a642a1htto8z",
	}

	if err := client.Do(token, retrieve); err != nil {
		if errDb := s.repository.Create(entity.Transactions{
			Amount:   req.Amount,
			Currency: req.Currency,
			Token:    token.ID,
			Status:   "UNSUCCESS",
		}); errDb != nil {
			return err
		}
		return err
	}

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   req.Amount,
		Currency: req.Currency,
		Card:     token.ID,
	}
	if err := client.Do(charge, createCharge); err != nil {
		if errDb := s.repository.Create(entity.Transactions{
			Amount:   req.Amount,
			Currency: req.Currency,
			Token:    token.ID,
			Status:   "UNSUCCESS",
		}); errDb != nil {
			return err
		}
		return err
	}

	if errDb := s.repository.Create(entity.Transactions{
		Amount:   req.Amount,
		Currency: req.Currency,
		Token:    token.ID,
		Status:   "SUCCESS",
	}); errDb != nil {
		return err
	}

	return nil
}

func (s service) Inquiry(status string) ([]entity.Transactions, error) {
	return s.repository.Inquiry(status)
}
