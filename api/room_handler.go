package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
	"net/http"
	"time"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Rooms.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	err := c.BodyParser(&params)
	if err != nil {
		return c.Status(400).JSON(map[string]string{
			"error": "Bad request",
		})
	}

	if err := params.validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "Unauthorized",
		})
	}

	isAvailable, err := h.isRoomAvailable(c, params, roomID)
	if err != nil {
		return err
	}
	if !isAvailable {
		return c.Status(404).JSON(map[string]string{
			"error": "Cannot book this room in these dates",
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Bookings.InsertBooking(c.Context(), &booking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{
			"error": "Internal server error",
		})
	}
	fmt.Println(inserted)
	return nil
}

func (h *RoomHandler) isRoomAvailable(c *fiber.Ctx, params BookRoomParams, roomID primitive.ObjectID) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Bookings.GetBookings(c.Context(), where)
	if err != nil {
		return false, c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": "Room was not found",
		})
	}

	return len(bookings) > 0, nil
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("invalid date range")
	}
	if p.NumPersons <= 0 {
		return fmt.Errorf("invalid number of persons")
	}
	return nil
}
