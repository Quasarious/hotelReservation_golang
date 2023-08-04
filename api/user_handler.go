package api

import (
	"github.com/gofiber/fiber/v2"
	"hotelReservation_golang/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Iskander",
		LastName:  "Miftakhutdinov",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Iskander")
}
