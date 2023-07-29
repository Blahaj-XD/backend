package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) KidDashboardBankAccountInfo(c *fiber.Ctx) error {
	kidID, err := c.ParamsInt("kidID")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
	}

	output, err := s.backend.KidDashboardBankAccountInfo(c.Context(), backend.KidDashboardBankAccountInfoInput{
		KidID:       kidID,
		ParentID:    int(c.Locals("userID").(float64)),
		AccessToken: c.Locals("bankAccessToken").(string),
	})
	if err != nil {
		if errors.Is(err, backend.ErrNotFoundBankAccount) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, backend.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(output)
}
