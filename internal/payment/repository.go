package payment

import (
	"database/sql"
	"fmt"
	"inception/internal/entity"
	"strings"

	"github.com/labstack/gommon/log"
)

type Repository interface {
	Create(tran entity.Transactions) error
	Inquiry(status string) ([]entity.Transactions, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (r repository) Create(tran entity.Transactions) error {
	template := "INSERT INTO transactions(amount, currency,token,status) VALUES(%d,'%s','%s','%s')"
	sqlStr := fmt.Sprintf(template, tran.Amount, tran.Currency, tran.Token, tran.Status)
	_, err := r.db.Exec(sqlStr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func (r repository) Inquiry(status string) ([]entity.Transactions, error) {
	result := []entity.Transactions{}
	rows, err := r.db.Query(fmt.Sprintf("SELECT * FROM transactions where status = '%s'", strings.ToUpper(status)))
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.Transactions
		err = rows.Scan(&r.Amount, &r.Currency, &r.Status, &r.Token)
		if err != nil {
			log.Error(err)
		}
		result = append(result, r)
	}

	return result, nil
}
