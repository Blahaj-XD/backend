package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

type ParentAdminBankDepositKidAccountBody struct {
	Amount int `json:"amount"`
}

func (s *Server) ParentAdminBankDepositKidAccount(c *fiber.Ctx) error {
	var body ParentAdminBankDepositKidAccountBody

	kidID, err := c.ParamsInt("kidID")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid kid id")
	}

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := body.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	err = s.backend.ParentAdminBankDepositKidAccount(c.Context(), backend.ParentAdminBankDepositKidAccountInput{
		ParentID:    int(c.Locals("userID").(float64)),
		KidID:       kidID,
		Amount:      body.Amount,
		AccessToken: c.Locals("bankAccessToken").(string),
	})

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

func (b ParentAdminBankDepositKidAccountBody) Validate() error {
	if b.Amount == 0 {
		return errors.New("amount is required")
	}

	return nil
}
