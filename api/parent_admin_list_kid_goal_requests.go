package api

import (
	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) ParentAdminListKidGoalRequests(c *fiber.Ctx) error {
	parentID := int(c.Locals("userID").(float64))

	output, err := s.backend.ParentAdminListKidGoalRequests(c.Context(), backend.ParentAdminListKidGoalRequestsInput{
		ParentID: parentID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
