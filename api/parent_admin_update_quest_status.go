package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/BlahajXD/backend/repo"
	"github.com/gofiber/fiber/v2"
)

// CreateParentAdminUpdateQuestStatusHandler creates a handler for POST /parent-admin/quests/:id/{status}
func (s *Server) CreateParentAdminUpdateQuestStatusHandler(status repo.QuestStatus) fiber.Handler {
	return func(c *fiber.Ctx) error {
		parentID := int(c.Locals("userID").(float64))
		questID, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid quest id")
		}

		err = s.backend.ParentAdminUpdateQuestStatus(c.Context(), parentID, questID, status)
		if err != nil {
			if errors.Is(err, backend.ErrQuestNotFound) {
				return fiber.NewError(fiber.StatusNotFound, err.Error())
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "successfully updated quest status",
		})
	}
}
