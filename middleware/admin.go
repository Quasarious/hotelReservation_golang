package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hotelReservation_golang/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not enough rights")
	}
	if !user.IsAdmin {
		return fmt.Errorf("not enough rights")
	}

	return c.Next()
}
