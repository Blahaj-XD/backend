package api

import "github.com/gofiber/fiber/v2"

func (s *Server) AuthVerify(c *fiber.Ctx) error {
	verified := s.backend.AuthVerify(c.Locals("accessToken").(string))
	if !verified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"verified": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"verified": true,
	})
}
