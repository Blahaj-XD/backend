package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type Goal struct {
	ID            int
	KidID         int
	AccountNumber string
	Title         string
	TargetAmount  float64
	Status        GoalStatus
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
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

	// GoalStatusCancelled is a status for goal that is canceled
	// by the kid. (3)
	GoalStatusCancelled
)

func (s GoalStatus) String() string {
	switch s {
	case GoalStatusOngoing:
		return "ongoing"
	case GoalStatusAchieved:
		return "achieved"
	case GoalStatusOverdue:
		return "overdue"
	case GoalStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

func (d *Dependency) SaveGoal(ctx context.Context, params Goal) (Goal, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"kid_id",
		"account_number",
		"title",
		"target_amount",
		"status",
		"start_date",
		"end_date",
		"created_at"}

	query := qb.Insert("goals").
		Columns(cols...).
		Values(
			params.KidID, params.AccountNumber, params.Title,
			params.TargetAmount, params.Status,
			params.StartDate, params.EndDate, params.CreatedAt).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return Goal{}, errors.Wrap(err, "repo.SaveGoal")
	}

	var id int
	err = d.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return Goal{}, errors.Wrap(err, "repo.SaveGoal")
	}

	var output Goal
	output.ID = id
	output.KidID = params.KidID
	output.AccountNumber = params.AccountNumber
	output.Title = params.Title
	output.TargetAmount = params.TargetAmount
	output.Status = params.Status
	output.StartDate = params.StartDate
	output.EndDate = params.EndDate
	output.CreatedAt = params.CreatedAt

	return output, nil
}

func (d *Dependency) ListGoals(ctx context.Context, kidID int) ([]Goal, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"id",
		"kid_id",
		"account_number",
		"title",
		"target_amount",
		"status",
		"start_date",
		"end_date",
		"created_at"}

	query := qb.Select(cols...).
		From("goals").
		Where(sq.Eq{"kid_id": kidID}).
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "repo.ListGoals")
	}

	rows, err := d.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repo.ListGoals")
	}
	defer rows.Close()

	var output []Goal
	for rows.Next() {
		var item Goal
		err := rows.Scan(
			&item.ID,
			&item.KidID,
			&item.AccountNumber,
			&item.Title,
			&item.TargetAmount,
			&item.Status,
			&item.StartDate,
			&item.EndDate,
			&item.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "repo.ListGoals")
		}

		output = append(output, item)
	}

	return output, nil
}

func (d *Dependency) FindGoal(ctx context.Context, kidID, goalID int) (Goal, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"id",
		"kid_id",
		"account_number",
		"title",
		"target_amount",
		"status",
		"start_date",
		"end_date",
		"created_at"}

	query := qb.Select(cols...).
		From("goals").
		Where(sq.Eq{"id": goalID, "kid_id": kidID})

	sql, args, err := query.ToSql()
	if err != nil {
		return Goal{}, errors.Wrap(err, "repo.FindGoal")
	}

	var output Goal
	err = d.db.QueryRow(ctx, sql, args...).Scan(
		&output.ID,
		&output.KidID,
		&output.AccountNumber,
		&output.Title,
		&output.TargetAmount,
		&output.Status,
		&output.StartDate,
		&output.EndDate,
		&output.CreatedAt)
	if err != nil {
		return Goal{}, errors.Wrap(err, "repo.FindGoal")
	}

	return output, nil
}

func (d *Dependency) UpdateGoalStatus(ctx context.Context, kidID, goalID int, status GoalStatus) error {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := qb.Update("goals").
		Set("status", status).
		Where(sq.Eq{"id": goalID, "kid_id": kidID})

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "repo.UpdateGoalStatus")
	}

	_, err = d.db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "repo.UpdateGoalStatus")
	}

	return nil
}

type KidBalanceRequest struct {
	ID                int
	KidID             int
	ParentID          int
	FromAccountNumber string
	ToAccountNumber   string
	Title             string
	Description       string
	Amount            float64
	Status            KidBalanceRequestStatus
	CreatedAt         time.Time
}

type KidBalanceRequestStatus int

const (
	// KidBalanceRequestStatusPending when the request is created
	// and waiting for approval from the parent
	KidBalanceRequestStatusPending KidBalanceRequestStatus = iota

	// KidBalanceRequestStatusApproved when the request is approved
	// by the parent
	KidBalanceRequestStatusApproved

	// KidBalanceRequestStatusRejected when the request is rejected
	// by the parent
	KidBalanceRequestStatusRejected
)

func (s KidBalanceRequestStatus) String() string {
	switch s {
	case KidBalanceRequestStatusPending:
		return "pending"
	case KidBalanceRequestStatusApproved:
		return "approved"
	case KidBalanceRequestStatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}

func (d *Dependency) NewKidBalanceRequest(ctx context.Context, params KidBalanceRequest) (KidBalanceRequest, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"kid_id",
		"parent_id",
		"from_account_number",
		"to_account_number",
		"title",
		"description",
		"amount",
		"status",
		"created_at"}

	query := qb.Insert("kid_balance_requests").
		Columns(cols...).
		Values(
			params.KidID, params.ParentID, params.Title,
			params.FromAccountNumber, params.ToAccountNumber, params.Description,
			params.Amount, params.Status, params.CreatedAt).
		Suffix("RETURNING \"id\"")

	sql, args, err := query.ToSql()
	if err != nil {
		return KidBalanceRequest{}, errors.Wrap(err, "repo.NewKidBalanceRequest")
	}

	var id int
	err = d.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return KidBalanceRequest{}, errors.Wrap(err, "repo.NewKidBalanceRequest")
	}

	var output KidBalanceRequest
	output.ID = id
	output.KidID = params.KidID
	output.ParentID = params.ParentID
	output.Title = params.Title
	output.Description = params.Description
	output.Amount = params.Amount
	output.Status = params.Status
	output.CreatedAt = params.CreatedAt

	return output, nil
}

func (d *Dependency) ListKidBalanceRequest(ctx context.Context, col string, value any) ([]KidBalanceRequest, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"id",
		"kid_id",
		"parent_id",
		"from_account_number",
		"to_account_number",
		"title",
		"description",
		"amount",
		"status",
		"created_at",
	}

	query := qb.Select(cols...).
		From("kid_balance_requests").
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "repo.ListKidBalanceRequest")
	}

	rows, err := d.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repo.ListKidBalanceRequest")
	}
	defer rows.Close()

	var output []KidBalanceRequest
	for rows.Next() {
		var item KidBalanceRequest
		err := rows.Scan(
			&item.ID,
			&item.KidID,
			&item.ParentID,
			&item.FromAccountNumber,
			&item.ToAccountNumber,
			&item.Title,
			&item.Description,
			&item.Amount,
			&item.Status,
			&item.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "repo.ListKidBalanceRequest")
		}

		output = append(output, item)
	}

	return output, nil
}

func (d *Dependency) FindKidBalanceRequest(ctx context.Context, id int) (KidBalanceRequest, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"id",
		"kid_id",
		"parent_id",
		"from_account_number",
		"to_account_number",
		"title",
		"description",
		"amount",
		"status",
		"created_at",
	}

	query := qb.Select(cols...).
		From("kid_balance_requests").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return KidBalanceRequest{}, errors.Wrap(err, "repo.FindKidBalanceRequest")
	}

	var output KidBalanceRequest
	err = d.db.QueryRow(ctx, sql, args...).Scan(
		&output.ID,
		&output.KidID,
		&output.ParentID,
		&output.FromAccountNumber,
		&output.ToAccountNumber,
		&output.Title,
		&output.Description,
		&output.Amount,
		&output.Status,
		&output.CreatedAt)
	if err != nil {
		return KidBalanceRequest{}, errors.Wrap(err, "repo.FindKidBalanceRequest")
	}

	return output, nil
}

func (d *Dependency) UpdateKidBalanceRequestStatus(ctx context.Context, id int, status KidBalanceRequestStatus) error {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := qb.Update("kid_balance_requests").
		Set("status", status).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "repo.UpdateKidBalanceRequestStatus")
	}

	_, err = d.db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "repo.UpdateKidBalanceRequestStatus")
	}

	return nil
}
