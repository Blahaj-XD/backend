package api

import (
	"time"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

type KidDashboardCreateGoalBody struct {
	Title        string    `json:"title"`
	TargetAmount float64   `json:"target_amount"`
	EndDate      time.Time `json:"end_date"`
}

func (s *Server) KidDashboardCreateGoal(c *fiber.Ctx) error {
	kidID, err := c.ParamsInt("kidID")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
	}

	var body KidDashboardCreateGoalBody

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := body.Validate(); err != nil {
		return err
	}

	goal, err := s.backend.KidDashboardCreateGoal(c.Context(), backend.KidDashboardCreateGoalInput{
		KidID:        kidID,
		ParentID:     int(c.Locals("userID").(float64)),
		Title:        body.Title,
		TargetAmount: body.TargetAmount,
		EndDate:      body.EndDate,
		AccessToken:  c.Locals("bankAccessToken").(string),
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(goal)
}

func (body KidDashboardCreateGoalBody) Validate() error {
	if body.Title == "" {
		return fiber.NewError(fiber.StatusBadRequest, "title is required")
	}

	if body.TargetAmount <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "target_amount must be greater than 0")
	}

	if body.EndDate.IsZero() {
		return fiber.NewError(fiber.StatusBadRequest, "end_date is required")
	}

	return nil
}
