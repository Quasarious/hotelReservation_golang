package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"hotelReservation_golang/db"
	"net/http"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBookings(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Bookings.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "unauthorized",
		})
	}
	if err := h.store.Bookings.UpdateBooking(c.Context(), booking.ID.String(), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(map[string]string{
		"status": "updated",
	})
}

// for admin only
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Bookings.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// for user
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Bookings.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "unauthorized",
		})
	}
	return c.JSON(booking)
}
