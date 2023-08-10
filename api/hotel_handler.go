package api

import (
	"github.com/gofiber/fiber/v2"
	"hotelReservation_golang/db"
)

type HotelHandler struct {
	roomStore  db.RoomStore
	hotelStore db.HotelStore
}

type HotelQueryParams struct {
	Rooms  bool
	Rating float64
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		roomStore:  rs,
		hotelStore: hs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams

	if err := c.QueryParser(&qparams); err != nil {
		return err
	}

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}
