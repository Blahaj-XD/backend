package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

type ParentAdminListQuestsItem struct {
	ID          int     `json:"id"`
	ParentID    int     `json:"parent_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Reward      float64 `json:"reward"`
	Status      struct {
		Code repo.QuestStatus `json:"code"`
		Name string           `json:"name"`
	} `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}

type ParentAdminListQuestsOutput struct {
	TotalItems int                         `json:"total_items"`
	Items      []ParentAdminListQuestsItem `json:"items"`
}

func (d *Dependency) ParentAdminListQuests(ctx context.Context, parentID int) (ParentAdminListQuestsOutput, error) {
	var output ParentAdminListQuestsOutput

	quests, err := d.repo.ListQuests(ctx, parentID)
	if err != nil {
		return ParentAdminListQuestsOutput{}, errors.Wrap(err, "backend.ParentAdminListQuests -> repo.ListQuests")
	}

	output.TotalItems = len(quests)
	output.Items = make([]ParentAdminListQuestsItem, 0)

	for _, quest := range quests {
		// Check if quest has past end date and update status to expired
		// This might be better if it's done in a cron job or other background process
		if quest.Status == repo.QuestStatusAvailable && quest.EndDate.Before(time.Now()) {
			quest.Status = repo.QuestStatusExpired
			err := d.repo.UpdateQuestStatus(ctx, parentID, quest.ID, quest.Status)
			if err != nil {
				return ParentAdminListQuestsOutput{}, errors.Wrap(err, "backend.ParentAdminListQuests -> repo.UpdateQuestStatus")
			}
		}

		var item ParentAdminListQuestsItem
		item.ID = quest.ID
		item.ParentID = quest.ParentID
		item.Title = quest.Title
		item.Description = quest.Description
		item.Reward = quest.Reward
		item.Status.Code = quest.Status
		item.Status.Name = quest.Status.String()
		item.StartDate = quest.StartDate
		item.EndDate = quest.EndDate
		item.CreatedAt = quest.CreatedAt

		output.Items = append(output.Items, item)
	}

	return output, nil
}
