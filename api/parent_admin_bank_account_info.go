package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) ParentAdminBankAccountInfo(c *fiber.Ctx) error {
	output, err := s.backend.ParentAdminBankAccountInfo(c.Context(), backend.ParentAdminBankAccountInfoInput{
		ParentID:    int(c.Locals("userID").(float64)),
		AccessToken: c.Locals("bankAccessToken").(string),
	})
	if err != nil {
		if errors.Is(err, backend.ErrNotFoundBankAccount) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(output)
}
