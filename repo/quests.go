package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type Quest struct {
	ID          int
	ParentID    int
	Title       string
	Description string
	Reward      float64
	Status      QuestStatus
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
}

// QuestStatus represents the status of a quest.
type QuestStatus int

const (
	// QuestStatusAvailable is a status for quest that is available to be (0)
	// taken by the kid.
	QuestStatusAvailable QuestStatus = iota

	// QuestStatusOngoing is a status for quest that is taken by the kid. (1)
	QuestStatusOngoing

	// QuestStatusDone marks the quest as done. (2)
	QuestStatusDone

	// QuestStatusApproved is a special status for quest that is approved by
	// the parent after the kid has done it. (3)
	QuestStatusApproved

	// QuestStatusCanceled is a special status for quest that is canceled
	// by the parent. (4)
	QuestStatusCanceled

	// QuestStatusExpired is a special status for quest that is expired. (5)
	QuestStatusExpired
)

func (s QuestStatus) String() string {
	switch s {
	case QuestStatusAvailable:
		return "available"
	case QuestStatusOngoing:
		return "ongoing"
	case QuestStatusDone:
		return "done"
	case QuestStatusApproved:
		return "approved"
	case QuestStatusCanceled:
		return "canceled"
	case QuestStatusExpired:
		return "expired"
	default:
		return ""
	}
}

func (d *Dependency) SaveQuest(ctx context.Context, params Quest) error {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols := []string{
		"parent_id",
		"title",
		"description",
		"reward",
		"status",
		"start_date",
		"end_date",
		"created_at"}

	sql, args, err := qb.Insert("quests").
		Columns(cols...).
		Values(params.ParentID, params.Title, params.Description, params.Reward, params.Status, params.StartDate, params.EndDate, params.CreatedAt).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "repo.SaveQuest: ToSql")
	}

	_, err = d.db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "repo.SaveQuest: exec sql")
	}

	return nil
}

func (d *Dependency) ListQuests(ctx context.Context, parentID int) ([]Quest, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := qb.Select("id", "title", "description", "reward", "status", "start_date", "end_date", "created_at").
		From("quests").
		Where(sq.Eq{"parent_id": parentID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "repo.ListQuests: ToSql")
	}

	rows, err := d.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repo.ListQuests: query sql")
	}
	defer rows.Close()

	var quests []Quest
	for rows.Next() {
		var quest Quest
		err := rows.Scan(&quest.ID, &quest.Title, &quest.Description, &quest.Reward, &quest.Status, &quest.StartDate, &quest.EndDate, &quest.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "repo.ListQuests: scan row")
		}
		quests = append(quests, quest)
	}

	return quests, nil
}

func (d *Dependency) FindQuest(ctx context.Context, parentID, questID int) (Quest, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := qb.Select("id", "title", "description", "reward", "status", "start_date", "end_date", "created_at").
		From("quests").
		Where(sq.And{
			sq.Eq{"id": questID},
			sq.Eq{"parent_id": parentID},
		}).
		ToSql()
	if err != nil {
		return Quest{}, errors.Wrap(err, "repo.FindQuest: ToSql")
	}

	row := d.db.QueryRow(ctx, sql, args...)
	var quest Quest
	err = row.Scan(&quest.ID, &quest.Title, &quest.Description, &quest.Reward, &quest.Status, &quest.StartDate, &quest.EndDate, &quest.CreatedAt)
	if err != nil {
		return Quest{}, errors.Wrap(err, "repo.FindQuest: scan row")
	}

	return quest, nil
}

func (d *Dependency) UpdateQuestStatus(ctx context.Context, parentID, questID int, status QuestStatus) error {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := qb.Update("quests").
		Set("status", status).
		Where(sq.And{
			sq.Eq{"id": questID},
			sq.Eq{"parent_id": parentID},
		}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "repo.UpdateQuestStatus: ToSql")
	}

	_, err = d.db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "repo.UpdateQuestStatus: exec sql")
	}

	return nil
}

func (d *Dependency) DeleteQuest(ctx context.Context, questID int) error {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := qb.Delete("quests").
		Where(sq.Eq{"id": questID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "repo.DeleteQuest: ToSql")
	}

	_, err = d.db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "repo.DeleteQuest: exec sql")
	}

	return nil
}
