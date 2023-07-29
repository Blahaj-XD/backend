package backend

import (
	"context"

	"github.com/pkg/errors"
)

type KidDashboardBankTransactionInfoInput struct {
	KidID          int `json:"kid_id"`
	ParentID       int `json:"parent_id"`
	Page           int `json:"page"`
	RecordsPerPage int `json:"records_per_page"`
	AccessToken    string
}

type KidDashboardBankTransactionInfoItem struct {
	UID                   int     `json:"uid"`
	Amount                float64 `json:"amount"`
	AccountNumber         string  `json:"account_number"`
	ReceiverAccountNumber string  `json:"receiver_account_number"`
	TransactionType       string  `json:"transaction_type"`
}

type KidDashboardBankTransactionInfoOutput struct {
	TotalItems int                                   `json:"total_items"`
	Items      []KidDashboardBankTransactionInfoItem `json:"items"`
}

func (d *Dependency) KidDashboardBankTransactionInfo(ctx context.Context, input KidDashboardBankTransactionInfoInput) (KidDashboardBankTransactionInfoOutput, error) {
	kid, err := d.repo.FindKid(ctx, "id", input.KidID)
	if err != nil {
		return KidDashboardBankTransactionInfoOutput{}, ErrUserNotFound
	}

	if kid.ParentID != input.ParentID {
		return KidDashboardBankTransactionInfoOutput{}, ErrUserNotFound
	}

	transactions, err := d.BankTransactionInfo(BankTransactionInfoInput{
		AccountNumber:  kid.AccountNumber,
		Page:           input.Page,
		RecordsPerPage: input.RecordsPerPage,
		AccessToken:    input.AccessToken,
	})
	if err != nil {
		return KidDashboardBankTransactionInfoOutput{}, errors.Wrap(err, "backend.KidDashboardBankTransactionInfo -> BankTransactionInfo")
	}

	var output KidDashboardBankTransactionInfoOutput
	output.TotalItems = transactions.TotalItems
	output.Items = make([]KidDashboardBankTransactionInfoItem, 0)

	for _, transaction := range transactions.Items {
		var item KidDashboardBankTransactionInfoItem
		item.UID = transaction.UID
		item.Amount = transaction.Amount
		item.TransactionType = transaction.TransactionType
		item.AccountNumber = transaction.AccountNo
		item.ReceiverAccountNumber = transaction.ReceiverAccountNo

		output.Items = append(output.Items, item)
	}

	return output, nil
}
