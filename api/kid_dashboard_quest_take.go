package api

import (
	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *Server) KidDashboardQuestTake(c *fiber.Ctx) error {
	parentID := int(c.Locals("userID").(float64))

	kidID, err := c.ParamsInt("kidID")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
	}

	questID, err := c.ParamsInt("questID")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
	}

	err = s.backend.KidDashboardQuestTake(c.Context(), backend.KidDashboardQuestTakeInput{
		ParentID: parentID,
		KidID:    kidID,
		QuestID:  questID,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to take quest")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "successfully took quest",
	})
}
