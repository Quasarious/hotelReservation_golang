package api

import (
	"github.com/gofiber/fiber/v2"
	"hotelReservation_golang/types"
	"net/http"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "unauthorized",
		})
	}

	return user, nil
}
