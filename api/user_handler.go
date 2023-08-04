package api

import (
	"github.com/gofiber/fiber/v2"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
)

type UserHandler struct {
	userStore db.UserStorage
}

func NewUserHandler(storage db.UserStorage) *UserHandler {
	return &UserHandler{
		userStore: storage,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)

	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)

	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserByID(c.Context(), id)

	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())

	if err != nil {
		return err
	}

	return c.JSON(users)
}
