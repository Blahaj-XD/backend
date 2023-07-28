package backend

import (
	"context"
	"database/sql"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

func (d *Dependency) ParentAdminUpdateQuestStatus(
	ctx context.Context, parentID, questID int, status repo.QuestStatus,
) error {
	if _, err := d.repo.FindQuest(ctx, parentID, questID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrQuestNotFound
		}
		return errors.Wrap(err, "backend.ParentAdminUpdateQuestStatus -> repo.FindQuest")
	}

	if err := d.repo.UpdateQuestStatus(ctx, parentID, questID, status); err != nil {
		return errors.Wrap(err, "backend.ParentAdminUpdateQuestStatus -> repo.UpdateQuestStatus")
	}

	return nil
}
