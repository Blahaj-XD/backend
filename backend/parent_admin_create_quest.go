package backend

import (
	"context"
	"time"

	"github.com/BlahajXD/backend/repo"
	"github.com/pkg/errors"
)

type ParentAdminCreateQuestInput struct {
	ParentID    int       `json:"parent_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Reward      float64   `json:"reward"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type ParentAdminCreateQuestOutput struct {
	ID          int              `json:"id"`
	ParentID    int              `json:"parent_id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Reward      float64          `json:"reward"`
	Status      repo.QuestStatus `json:"status"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	CreatedAt   time.Time        `json:"created_at"`
}

func (d *Dependency) ParentAdminCreateQuest(ctx context.Context, input ParentAdminCreateQuestInput) (ParentAdminCreateQuestOutput, error) {
	params := repo.Quest{
		ParentID:    input.ParentID,
		Title:       input.Title,
		Description: input.Description,
		Reward:      input.Reward,
		Status:      repo.QuestStatusAvailable,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		CreatedAt:   time.Now(),
	}

	err := d.repo.SaveQuest(ctx, params)
	if err != nil {
		return ParentAdminCreateQuestOutput{}, errors.Wrap(err, "backend.ParentAdminCreateQuest -> repo.SaveQuest")
	}

	output := ParentAdminCreateQuestOutput{
		ID:          params.ID,
		ParentID:    params.ParentID,
		Title:       params.Title,
		Description: params.Description,
		Reward:      params.Reward,
		Status:      params.Status,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		CreatedAt:   params.CreatedAt,
	}

	return output, nil
}
