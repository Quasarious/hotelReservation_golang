package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
)

type UserHandler struct {
	store db.Store
}

func NewUserHandler(storage *db.Store) *UserHandler {
	return &UserHandler{
		store: *storage,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	if err = c.BodyParser(&params); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	if err = h.store.Users.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated: ": userID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.store.Users.DeleteUser(c.Context(), userID); err != nil {
		return err
	}

	return c.JSON(map[string]string{"deleted:": userID})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errs := params.Validate(); len(errs) > 0 {
		return c.JSON(errs)
	}

	user, err := types.NewUserFromParams(params)

	if err != nil {
		return err
	}

	insertedUser, err := h.store.Users.InsertUser(c.Context(), user)

	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.store.Users.GetUserByID(c.Context(), id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error:": "not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.Users.GetUsers(c.Context())

	if err != nil {
		return err
	}

	return c.JSON(users)
}
