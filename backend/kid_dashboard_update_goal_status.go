package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

type KidDashboardUpdateGoalStatusInput struct {
	ParentID int             `json:"parent_id"`
	KidID    int             `json:"kid_id"`
	GoalID   int             `json:"goal_id"`
	Status   repo.GoalStatus `json:"status"`
}

func (d *Dependency) KidDashboardUpdateGoalStatus(ctx context.Context, input KidDashboardUpdateGoalStatusInput) error {
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

	// Force overdue status if the goal is overdue
	if goal.EndDate.Before(time.Now()) {
		input.Status = repo.GoalStatusOverdue
	}

	if err := d.repo.UpdateGoalStatus(ctx, kid.ID, goal.ID, input.Status); err != nil {
		return errors.Wrap(err, "backend.KidDashboardUpdateGoalStatus -> repo.UpdateGoalStatus")
	}

	return nil
}
