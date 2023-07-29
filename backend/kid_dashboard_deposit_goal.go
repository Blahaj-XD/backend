package backend

import (
	"context"

	"github.com/BlahajXD/backend/repo"
)

type KidDashboardDepositGoalInput struct {
	ParentID    int `json:"parent_id"`
	KidID       int `json:"kid_id"`
	GoalID      int `json:"goal_id"`
	Amount      int `json:"amount"`
	AccessToken string
}

func (d *Dependency) KidDashboardDepositGoal(ctx context.Context, input KidDashboardDepositGoalInput) error {
	kid, err := d.repo.FindKid(ctx, "id", input.KidID)
	if err != nil {
		return ErrUserNotFound
	}

	if kid.ParentID != input.ParentID {
		return ErrUserNotFound
	}

	goal, err := d.repo.FindGoal(ctx, input.KidID, input.GoalID)
	if err != nil {
		return ErrGoalNotFound
	}

	if goal.Status != repo.GoalStatusOngoing {
		return ErrGoalNotFound
	}

	err = d.BankAccountCreateTransaction(ctx, BankAccountCreateTransactionInput{
		FromAccountNumber: kid.AccountNumber,
		ToAccountNumber:   goal.AccountNumber,
		Amount:            input.Amount,
		AccessToken:       input.AccessToken,
	})
	if err != nil {
		return err
	}

	return nil
}
