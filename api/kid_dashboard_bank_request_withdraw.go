package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

type KidDashboardBankRequestWithdrawBody struct {
	FromAccountNumber string `json:"from_account_number"`
	ToAccountNumber   string `json:"account_number"`
}

func (s *Server) KidDashboardBankRequestWithdraw(c *fiber.Ctx) error {
	kidID, err := c.ParamsInt("kidID")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
	}
	var body KidDashboardBankRequestWithdrawBody

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := body.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	err = s.backend.KidDashboardBankRequestWithdraw(c.Context(), backend.KidDashboardBankRequestWithdrawInput{
		ParentID:          int(c.Locals("userID").(float64)),
		KidID:             kidID,
		FromAccountNumber: body.FromAccountNumber,
		ToAccountNumber:   body.ToAccountNumber,
		AccessToken:       c.Locals("bankAccessToken").(string),
	})
	if err != nil {
		if errors.Is(err, backend.ErrGoalNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "successfully requested goal withdraw",
	})
}

func (b KidDashboardBankRequestWithdrawBody) Validate() error {
	if b.FromAccountNumber == "" {
		return errors.New("from account number is required")
	}

	if b.ToAccountNumber == "" {
		return errors.New("to account number is required")
	}

	return nil
}
