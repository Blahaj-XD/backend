package api

import (
	"errors"

	"github.com/BlahajXD/backend/backend"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) ParentAdminBankTransactionInfo(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	recordsPerPage := c.QueryInt("records_per_page", 10)

	output, err := s.backend.ParentAdminBankTransactionInfo(c.Context(),
		backend.ParentAdminBankTransactionInfoInput{
			ParentID:       int(c.Locals("userID").(float64)),
			Page:           page,
			RecordsPerPage: recordsPerPage,
			AccessToken:    c.Locals("bankAccessToken").(string),
		})
	if err != nil {
		if errors.Is(err, backend.ErrNotFoundBankAccount) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, backend.ErrBankTransactionNotExists) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(output)
}
