package repo

import "time"

type Goal struct {
	ID           int
	KidID        int
	Title        string
	Description  string
	TargetAmount float64
	StartDate    time.Time
	EndDate      time.Time
	CreatedAt    time.Time
}

// GoalStatus represents the status of a goal.
type GoalStatus int

const (
	// GoalStatusOngoing is a status for goal that is ongoing. (0)
	GoalStatusOngoing GoalStatus = iota

	// GoalStatusAchieved is a status for goal that has
	// achieved the target amount. (1)
	GoalStatusAchieved

	// GoalStatusOverdue is a status for goal that is overdue
	// (the end date has passed). (2)
	GoalStatusOverdue

	// GoalStatusCanceled is a status for goal that is canceled
	// by the kid. (3)
	GoalStatusCanceled
)
