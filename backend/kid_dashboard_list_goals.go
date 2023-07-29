package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

type KidDashboardListGoalsItem struct {
	ID            int       `json:"id"`
	KidID         int       `json:"kid_id"`
	AccountNumber string    `json:"account_number"`
	Title         string    `json:"title"`
	Current       float64   `json:"current"`
	TargetAmount  float64   `json:"target_amount"`
	Status        string    `json:"status"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CreatedAt     time.Time `json:"created_at"`
}

type KidDashboardListGoalsOutput struct {
	TotalItems int                         `json:"total_items"`
	Items      []KidDashboardListGoalsItem `json:"items"`
}

type KidDashboardListGoalsInput struct {
	ParentID    int `json:"parent_id"`
	KidID       int `json:"kid_id"`
	AccessToken string
}

func (d *Dependency) KidDashboardListGoals(ctx context.Context, input KidDashboardListGoalsInput) (KidDashboardListGoalsOutput, error) {
	kid, err := d.repo.FindKid(ctx, "id", input.KidID)
	if err != nil {
		return KidDashboardListGoalsOutput{}, ErrUserNotFound
	}

	if kid.ParentID != input.ParentID {
		return KidDashboardListGoalsOutput{}, ErrUserNotFound
	}

	goals, err := d.repo.ListGoals(ctx, input.KidID)
	if err != nil {
		return KidDashboardListGoalsOutput{}, err
	}

	var output KidDashboardListGoalsOutput
	output.TotalItems = len(goals)
	output.Items = make([]KidDashboardListGoalsItem, 0)

	for _, goal := range goals {
		info, err := d.BankAccountInfo(ctx, BankAccountInfoInput{
			AccessToken:   input.AccessToken,
			AccountNumber: goal.AccountNumber,
		})
		if err != nil {
			return KidDashboardListGoalsOutput{}, err
		}

		// check if goal is overdue
		if goal.Status == repo.GoalStatusOngoing && time.Now().After(goal.EndDate) {
			goal.Status = repo.GoalStatusOverdue
			err := d.repo.UpdateGoalStatus(ctx, kid.ID, goal.ID, repo.GoalStatusOverdue)
			if err != nil {
				return KidDashboardListGoalsOutput{}, errors.Wrap(err, "KidDashboardListGoals: update goal status")
			}
		}

		// mark goal as completed if target amount is reached before end date
		if goal.Status == repo.GoalStatusOngoing && info.Balance >= goal.TargetAmount {
			goal.Status = repo.GoalStatusAchieved
			err := d.repo.UpdateGoalStatus(ctx, kid.ID, goal.ID, repo.GoalStatusAchieved)
			if err != nil {
				return KidDashboardListGoalsOutput{}, errors.Wrap(err, "KidDashboardListGoals: update goal status")
			}
		}

		output.Items = append(output.Items, KidDashboardListGoalsItem{
			ID:            goal.ID,
			KidID:         goal.KidID,
			AccountNumber: goal.AccountNumber,
			Title:         goal.Title,
			Current:       info.Balance,
			TargetAmount:  goal.TargetAmount,
			Status:        goal.Status.String(),
			StartDate:     goal.StartDate,
			EndDate:       goal.EndDate,
			CreatedAt:     goal.CreatedAt,
		})
	}

	return output, nil
}
