package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

type KidDashboardCreateGoalInput struct {
	KidID        int       `json:"kid_id"`
	ParentID     int       `json:"parent_id"`
	Title        string    `json:"title"`
	TargetAmount float64   `json:"target_amount"`
	EndDate      time.Time `json:"end_date"`
	AccessToken  string
}

type KidDashboardCreateGoalOutput struct {
	ID            int       `json:"id"`
	KidID         int       `json:"kid_id"`
	AccountNumber string    `json:"account_number"`
	Title         string    `json:"title"`
	TargetAmount  float64   `json:"target_amount"`
	Status        string    `json:"status"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CreatedAt     time.Time `json:"created_at"`
}

func (d *Dependency) KidDashboardCreateGoal(ctx context.Context, input KidDashboardCreateGoalInput) (KidDashboardCreateGoalOutput, error) {
	kid, err := d.repo.FindKid(ctx, "id", input.KidID)
	if err != nil {
		return KidDashboardCreateGoalOutput{}, ErrUserNotFound
	}

	if kid.ParentID != input.ParentID {
		return KidDashboardCreateGoalOutput{}, ErrUserNotFound
	}

	accountNumber, err := d.BankCreateBankAccount(BankCreateBankAccountInput{
		Balance:     0,
		AccessToken: input.AccessToken,
	})
	if err != nil {
		return KidDashboardCreateGoalOutput{}, errors.Wrap(err, "KidDashboardCreateGoal: create bank account")
	}

	now := time.Now()
	goal, err := d.repo.SaveGoal(ctx, repo.Goal{
		KidID:         input.KidID,
		AccountNumber: accountNumber,
		Title:         input.Title,
		TargetAmount:  input.TargetAmount,
		Status:        repo.GoalStatusOngoing,
		StartDate:     now,
		EndDate:       input.EndDate,
		CreatedAt:     now,
	})
	if err != nil {
		return KidDashboardCreateGoalOutput{}, errors.Wrap(err, "KidDashboardCreateGoal: save goal")
	}

	return KidDashboardCreateGoalOutput{
		ID:            goal.ID,
		KidID:         goal.KidID,
		AccountNumber: goal.AccountNumber,
		Title:         goal.Title,
		TargetAmount:  goal.TargetAmount,
		Status:        goal.Status.String(),
		StartDate:     goal.StartDate,
		EndDate:       goal.EndDate,
		CreatedAt:     goal.CreatedAt,
	}, nil
}
