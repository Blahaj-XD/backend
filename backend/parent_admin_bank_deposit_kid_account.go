package backend

import (
	"context"

	"github.com/pkg/errors"
)

type ParentAdminBankDepositKidAccountInput struct {
	ParentID    int `json:"parent_id"`
	KidID       int `json:"kid_id"`
	Amount      int `json:"amount"`
	AccessToken string
}

func (d *Dependency) ParentAdminBankDepositKidAccount(ctx context.Context, input ParentAdminBankDepositKidAccountInput) error {
	parent, err := d.repo.FindParent(ctx, "id", input.ParentID)
	if err != nil {
		return errors.Wrap(ErrUserNotFound, "parent_admin_deposit_kid_account: failed to find parent")
	}

	kid, err := d.repo.FindKid(ctx, "id", input.KidID)
	if err != nil {
		return errors.Wrap(ErrUserNotFound, "parent_admin_deposit_kid_account: failed to find kid")
	}

	if kid.ParentID != parent.ID {
		return errors.Wrap(ErrUserNotFound, "parent_admin_deposit_kid_account: parent is not the owner of the kid")
	}

	err = d.BankAccountCreateTransaction(ctx, BankAccountCreateTransactionInput{
		FromAccountNumber: parent.AccountNumber,
		ToAccountNumber:   kid.AccountNumber,
		Amount:            input.Amount,
		AccessToken:       input.AccessToken,
	})
	if err != nil {
		return err
	}

	return nil
}
