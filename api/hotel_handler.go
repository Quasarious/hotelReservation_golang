package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservation_golang/db"
	"net/http"
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
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "no filtration able"})
	}

	filter := bson.M{}
	if params.Rating != 0.0 {
		filter["rating"] = params.Rating
	}

	hotels, err := h.store.Hotels.GetHotels(c.Context(), filter, &params.PaginationFilter)
	fmt.Printf("%+v \n", hotels[0])
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "No hotels found"})
	}

	resp := HotelsSourceResp{
		ResourceResp: ResourceResp{
			Results: len(hotels),
			Page:    int(params.Page),
		},
		Hotels: hotels,
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	var params RoomQueryParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "no filtration able"})
	}

	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Rooms.GetRooms(c.Context(), filter, &params.PaginationFilter)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "No rooms found"})
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid ID"})
	}

	hotel, err := h.store.Hotels.GetHotelByID(c.Context(), oid)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "No hotels found with such ID"})
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	err := h.store.Hotels.DeleteHotel(c.Context(), hotelID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "No hotels found with such ID"})
	}

	return c.JSON(map[string]string{"deleted:": "Hotel with ID: " + hotelID})
}
