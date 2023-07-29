package middleware

import (
	"strings"

	"github.com/BlahajXD/backend/logic"
	"github.com/gofiber/fiber/v2"
)

func Authenticated(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if len(strings.Split(authorization, " ")) < 2 {
		return fiber.NewError(fiber.StatusUnauthorized, "Please login to continue")
	}
	accessToken := strings.Split(authorization, " ")[1]

	bankAuthorization := c.Get("X-Bank-Authorization")
	if len(strings.Split(bankAuthorization, " ")) < 2 {
		return fiber.NewError(fiber.StatusUnauthorized, "Please login to continue")
	}

	bankAccessToken := strings.Split(bankAuthorization, " ")[1]

	_, claims, err := logic.VerifyJWT(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Please login to continue")
	}

	userID, ok := claims["userID"]
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Please try login to continue")
	}

	email, ok := claims["email"]
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Please try login to continue")
	}

	c.Locals("userID", userID)
	c.Locals("email", email)
	c.Locals("accessToken", accessToken)
	c.Locals("bankAccessToken", bankAccessToken)

	return c.Next()
}
