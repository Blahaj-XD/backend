package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/BlahajXD/backend/repo"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) CreateKidDashboardUpdateGoalStatusHandler(status repo.GoalStatus) fiber.Handler {
	return func(c *fiber.Ctx) error {
		kidID, err := c.ParamsInt("kidID")
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
		}

		goalID, err := c.ParamsInt("goalID")
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid goal id")
		}

		err = s.backend.KidDashboardUpdateGoalStatus(c.Context(), backend.KidDashboardUpdateGoalStatusInput{
			ParentID: int(c.Locals("userID").(float64)),
			KidID:    kidID,
			GoalID:   goalID,
			Status:   status,
		})
		if err != nil {
			if errors.Is(err, backend.ErrGoalNotFound) {
				return fiber.NewError(fiber.StatusNotFound, err.Error())
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "successfully updated goal status",
		})
	}
}
