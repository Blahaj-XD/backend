package api

import (
	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) KidDashboardListGoals(c *fiber.Ctx) error {
	parentID := int(c.Locals("userID").(float64))
	kidID, err := c.ParamsInt("kidID")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
	}

	output, err := s.backend.KidDashboardListGoals(c.Context(), backend.KidDashboardListGoalsInput{
		ParentID:    parentID,
		KidID:       kidID,
		AccessToken: c.Locals("bankAccessToken").(string),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
