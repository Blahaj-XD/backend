package api

import "github.com/gofiber/fiber/v2"

func (s *Server) ParentAdminListKids(c *fiber.Ctx) error {
	parentID := int(c.Locals("userID").(float64))

	output, err := s.backend.ParentAdminListKids(c.Context(), parentID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
