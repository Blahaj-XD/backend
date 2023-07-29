package backend

import (
	"context"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

type ParentAdminDecideKidGoalRequestInput struct {
	ParentID    int                          `json:"parent_id"`
	RequestID   int                          `json:"request_id"`
	Decision    repo.KidBalanceRequestStatus `json:"decision"`
	AccessToken string
}

func (d *Dependency) ParentAdminDecideKidGoalRequest(ctx context.Context, input ParentAdminDecideKidGoalRequestInput) error {
	parent, err := d.repo.FindParent(ctx, "id", input.ParentID)
	if err != nil {
		return ErrUserNotFound
	}

	request, err := d.repo.FindKidBalanceRequest(ctx, input.RequestID)
	if err != nil {
		return err
	}

	if request.ParentID != parent.ID {
		return errors.Wrap(ErrUserNotFound, "parent_admin_decide_kid_goal_request: parent is not the owner of the kid")
	}

	err = d.repo.UpdateKidBalanceRequestStatus(ctx, input.RequestID, input.Decision)
	if err != nil {
		return errors.Wrap(err, "ParentAdminDecideKidGoalRequest")
	}

	if input.Decision == repo.KidBalanceRequestStatusRejected {
		return nil
	}

	// parent approves the withdrawal request
	err = d.BankAccountCreateTransaction(ctx, BankAccountCreateTransactionInput{
		FromAccountNumber: request.FromAccountNumber,
		ToAccountNumber:   request.ToAccountNumber,
		Amount:            int(request.Amount),
		AccessToken:       input.AccessToken,
	})
	if err != nil {
		return err
	}

	return nil
}
