package backend

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ParentAdminBankTransactionInfoInput struct {
	ParentID       int `json:"parent_id"`
	Page           int `json:"page"`
	RecordsPerPage int `json:"records_per_page"`
	AccessToken    string
}

type ParentAdminBankTransactionInfoItem struct {
	UID                   int     `json:"uid"`
	Amount                float64 `json:"amount"`
	AccountNumber         string  `json:"account_number"`
	ReceiverAccountNumber string  `json:"receiver_account_number"`
	TransactionType       string  `json:"transaction_type"`
	KidName               string  `json:"kid_name"`
}

type ParentAdminBankTransactionInfoOutput struct {
	TotalItems int                                  `json:"total_items"`
	Items      []ParentAdminBankTransactionInfoItem `json:"items"`
}

func (d *Dependency) ParentAdminBankTransactionInfo(ctx context.Context, input ParentAdminBankTransactionInfoInput) (ParentAdminBankTransactionInfoOutput, error) {
	parent, err := d.repo.FindParent(ctx, "id", input.ParentID)
	if err != nil {
		return ParentAdminBankTransactionInfoOutput{}, ErrUserNotFound
	}

	transactions, err := d.BankTransactionInfo(BankTransactionInfoInput{
		AccountNumber:  parent.AccountNumber,
		Page:           input.Page,
		RecordsPerPage: input.RecordsPerPage,
		AccessToken:    input.AccessToken,
	})
	if err != nil {
		return ParentAdminBankTransactionInfoOutput{}, errors.Wrap(err, "backend.ParentAdminBankTransactionInfo -> BankTransactionInfo")
	}

	var output ParentAdminBankTransactionInfoOutput
	output.TotalItems = transactions.TotalItems
	output.Items = make([]ParentAdminBankTransactionInfoItem, 0)

	for _, transaction := range transactions.Items {
		var item ParentAdminBankTransactionInfoItem
		item.UID = transaction.UID
		item.Amount = transaction.Amount
		item.TransactionType = transaction.TransactionType

		// 0000000000 is the account number for top-ups
		if transaction.AccountNo == "0000000000" {
			item.KidName = "-"
		} else {
			log.Debug().Interface("transaction", transaction).Msg("backend.ParentAdminBankTransactionInfo")
			kid, err := d.repo.FindKid(ctx, "account_number", transaction.ReceiverAccountNo)
			if err != nil {
				log.Error().Err(err).Str("account_number", transaction.ReceiverAccountNo).Msg("backend.ParentAdminBankTransactionInfo -> FindKid")
				return ParentAdminBankTransactionInfoOutput{}, errors.Wrap(err, "backend.ParentAdminBankTransactionInfo -> FindKid")
			}
			item.KidName = kid.FullName
		}

		item.AccountNumber = transaction.AccountNo
		item.ReceiverAccountNumber = transaction.ReceiverAccountNo

		output.Items = append(output.Items, item)
	}

	return output, nil
}
