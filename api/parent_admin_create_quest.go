package api

import (
	"errors"
	"time"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

type ParentAdminCreateQuestBody struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Reward      float64   `json:"reward"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

func (s *Server) ParentAdminCreateQuest(c *fiber.Ctx) error {
	var body ParentAdminCreateQuestBody

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := body.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	parentID := int(c.Locals("userID").(float64))

	output, err := s.backend.ParentAdminCreateQuest(c.Context(), backend.ParentAdminCreateQuestInput{
		ParentID:    parentID,
		Title:       body.Title,
		Description: body.Description,
		Reward:      body.Reward,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

func (b ParentAdminCreateQuestBody) Validate() error {
	if b.Title == "" {
		return errors.New("title is required")
	}

	if b.Description == "" {
		return errors.New("description is required")
	}

	if b.Reward == 0 {
		return errors.New("reward is required")
	}

	if b.StartDate.IsZero() {
		return errors.New("start date is required")
	}

	if b.EndDate.IsZero() {
		return errors.New("end date is required")
	}

	return nil
}
