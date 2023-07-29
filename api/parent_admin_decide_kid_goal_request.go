package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/BlahajXD/backend/repo"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) CreateParentAdminDecideKidGoalRequestHandler(decision repo.KidBalanceRequestStatus) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := s.backend.ParentAdminDecideKidGoalRequest(c.Context(),
			backend.ParentAdminDecideKidGoalRequestInput{
				ParentID: c.Locals("user_id").(int),
			},
		)

		switch {
		case errors.Is(err, backend.ErrBankRequestParameter):
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, backend.ErrUserNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, backend.ErrNotFoundBankAccount):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, backend.ErrBankAmountNotEnough):
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		case err != nil:
			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "success",
		})
	}
}
