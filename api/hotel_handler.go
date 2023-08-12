package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservation_golang/db"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotels.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Rooms.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	hotel, err := h.store.Hotels.GetHotelByID(c.Context(), oid)
	if err != nil {
		return err
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	err := h.store.Hotels.DeleteHotel(c.Context(), hotelID)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"deleted:": c.Params("name") + " with ID: " + hotelID})
}
