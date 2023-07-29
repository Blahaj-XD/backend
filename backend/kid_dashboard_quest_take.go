package backend

import (
	"context"

	"github.com/BlahajXD/backend/repo"
)

type KidDashboardQuestTakeInput struct {
	ParentID int `json:"parent_id"`
	KidID    int `json:"kid_id"`
	QuestID  int `json:"quest_id"`
}

func (d *Dependency) KidDashboardQuestTake(ctx context.Context, input KidDashboardQuestTakeInput) error {
	quest, err := d.repo.FindQuest(ctx, input.ParentID, input.QuestID)
	if err != nil {
		return ErrQuestNotFound
	}

	if quest.Status != repo.QuestStatusAvailable {
		return ErrQuestNotAvailable
	}

	if err := d.repo.AssignKidToQuest(ctx, input.KidID, input.QuestID); err != nil {
		return err
	}

	return nil
}
