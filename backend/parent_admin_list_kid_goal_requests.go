package backend

import (
	"context"
	"time"
)

type ParentAdminListKidGoalRequestsInput struct {
	ParentID int `json:"parent_id"`
}

type ParentAdminListKidGoalRequestsItem struct {
	ID                int       `json:"id"`
	KidID             int       `json:"kid_id"`
	ParentID          int       `json:"parent_id"`
	FromAccountNumber string    `json:"from_account_number"`
	ToAccountNumber   string    `json:"to_account_number"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Amount            float64   `json:"amount"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type ParentAdminListKidGoalRequestsOutput struct {
	TotalItems int                                  `json:"total_items"`
	Items      []ParentAdminListKidGoalRequestsItem `json:"items"`
}

func (d *Dependency) ParentAdminListKidGoalRequests(ctx context.Context, input ParentAdminListKidGoalRequestsInput) (ParentAdminListKidGoalRequestsOutput, error) {
	var output ParentAdminListKidGoalRequestsOutput

	requests, err := d.repo.ListKidBalanceRequest(ctx, "parent_id", input.ParentID)
	if err != nil {
		return ParentAdminListKidGoalRequestsOutput{}, err
	}

	output.TotalItems = len(requests)
	output.Items = make([]ParentAdminListKidGoalRequestsItem, 0)
	for _, request := range requests {
		var item ParentAdminListKidGoalRequestsItem
		item.ID = request.ID
		item.KidID = request.KidID
		item.ParentID = request.ParentID
		item.FromAccountNumber = request.FromAccountNumber
		item.ToAccountNumber = request.ToAccountNumber
		item.Title = request.Title
		item.Description = request.Description
		item.Amount = request.Amount
		item.Status = request.Status.String()
		item.CreatedAt = request.CreatedAt

		output.Items = append(output.Items, item)
	}

	return output, nil
}
