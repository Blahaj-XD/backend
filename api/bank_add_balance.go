package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

type BankAddBalanceBody struct {
	Balance int `json:"balance"`
}

func (s *Server) BankAddBalance(c *fiber.Ctx) error {
	var body BankAddBalanceBody

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := body.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := s.backend.BankAddBalance(c.Context(), backend.BankAddBalanceInput{
		ParentID:    int(c.Locals("userID").(float64)),
		Balance:     body.Balance,
		AccessToken: c.Locals("bankAccessToken").(string),
	}); err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
	})
}

func (b BankAddBalanceBody) Validate() error {
	if b.Balance == 0 {
		return errors.New("balance is required")
	}
	return nil
}
