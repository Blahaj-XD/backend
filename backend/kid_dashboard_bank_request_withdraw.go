package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
)

type KidDashboardBankRequestWithdrawInput struct {
	ParentID          int    `json:"parent_id"`
	KidID             int    `json:"kid_id"`
	FromAccountNumber string `json:"account_number"`
	ToAccountNumber   string `json:"to_account_number"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	AccessToken       string
}

func (d *Dependency) KidDashboardBankRequestWithdraw(ctx context.Context, input KidDashboardBankRequestWithdrawInput) error {
	kid, err := d.repo.FindKid(ctx, "id", input.KidID)
	if err != nil {
		return ErrUserNotFound
	}

	if kid.ParentID != input.ParentID {
		return ErrUserNotFound
	}

	info, err := d.BankAccountInfo(ctx, BankAccountInfoInput{
		AccessToken:   input.AccessToken,
		AccountNumber: input.FromAccountNumber,
	})
	if err != nil {
		return err
	}

	params := repo.KidBalanceRequest{
		KidID:             kid.ID,
		ParentID:          kid.ParentID,
		FromAccountNumber: input.FromAccountNumber,
		ToAccountNumber:   input.ToAccountNumber,
		Title:             input.Title,
		Description:       input.Description,
		Amount:            info.Balance,
		Status:            repo.KidBalanceRequestStatusPending,
		CreatedAt:         time.Now(),
	}
	_, err = d.repo.NewKidBalanceRequest(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
