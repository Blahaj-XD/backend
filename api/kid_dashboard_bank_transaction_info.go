package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) KidDashboardBankTransactionInfo(c *fiber.Ctx) error {
	kidID, err := c.ParamsInt("kidID")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid kid id")
	}

	page := c.QueryInt("page", 1)
	recordsPerPage := c.QueryInt("records_per_page", 10)

	output, err := s.backend.KidDashboardBankTransactionInfo(c.Context(),
		backend.KidDashboardBankTransactionInfoInput{
			ParentID:       int(c.Locals("userID").(float64)),
			KidID:          kidID,
			Page:           page,
			RecordsPerPage: recordsPerPage,
			AccessToken:    c.Locals("bankAccessToken").(string),
		})
	if err != nil {
		if errors.Is(err, backend.ErrNotFoundBankAccount) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(output)
}
