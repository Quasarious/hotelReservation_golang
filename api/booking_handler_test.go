package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"hotelReservation_golang/middleware"
	"hotelReservation_golang/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBookingHandler_HandleGetBooking_CorrectID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user  = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room  = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		route          = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	route.Get("/", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var bookingResp *types.Booking

	err = json.NewDecoder(resp.Body).Decode(&bookingResp)
	assert.NoError(t, err)
	assert.Equal(t, booking.ID, bookingResp.ID)
	assert.Equal(t, booking.UserID, bookingResp.UserID)
	assert.Equal(t, booking.RoomID, bookingResp.RoomID)
}

func TestBookingHandler_HandleGetBooking_UnauthorizedUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		unauthorizedUser = AddUser(tdb.Store, "Jake", "Smith", false)
		user             = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel            = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room             = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(unauthorizedUser))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var bookingResp *types.Booking

	err = json.NewDecoder(resp.Body).Decode(&bookingResp)
	assert.NoError(t, err)
	assert.Empty(t, bookingResp)
}

func TestBookingHandler_HandleGetBooking_EmptyID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user  = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room  = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", ""), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var bookingResp *types.Booking

	err = json.NewDecoder(resp.Body).Decode(&bookingResp)
	assert.Error(t, err)
	assert.Empty(t, bookingResp)
}

func TestBookingHandler_HandleGetBookingsWithAccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser = AddUser(tdb.Store, "admin", "admin", true)
		user      = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel     = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room      = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var bookings []*types.Booking
	err = json.NewDecoder(resp.Body).Decode(&bookings)
	assert.NoError(t, err)
	assert.Equal(t, len(bookings), 1)
	assert.Equal(t, booking.ID, bookings[0].ID)
	assert.Equal(t, booking.UserID, bookings[0].UserID)
	assert.Equal(t, booking.RoomID, bookings[0].RoomID)

}

func TestBookingHandler_HandleGetBookingsWithoutAccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		anotherUser = AddUser(tdb.Store, "notadmin", "notadmin", false)
		user        = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel       = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room        = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(anotherUser))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var bookings []*types.Booking
	err = json.NewDecoder(resp.Body).Decode(&bookings)
	assert.Error(t, err)
	assert.Equal(t, len(bookings), 0)
}

func TestBookingHandler_HandleCancelBookings_CorrectData(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user  = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room  = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	route.Put("/:id", bookingHandler.HandleCancelBookings)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check if the booking was actually canceled in the database
	updatedBooking, err := tdb.Store.Bookings.GetBookingByID(context.Background(), booking.ID.Hex())
	assert.NoError(t, err)
	assert.True(t, updatedBooking.Canceled)
}

func TestBookingHandler_HandleCancelBookings_BookingDoesNotExist(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user  = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room  = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	route.Put("/:id", bookingHandler.HandleCancelBookings)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/%s", "fdsfsdf123sda"), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	var msg map[string]string
	err = json.NewDecoder(resp.Body).Decode(&msg)
	assert.NoError(t, err)
	assert.Equal(t, msg["error"], "Unable to find booking")

	// Check if the booking was actually canceled in the database
	updatedBooking, err := tdb.Store.Bookings.GetBookingByID(context.Background(), booking.ID.Hex())
	assert.NoError(t, err)
	assert.False(t, updatedBooking.Canceled)
}

func TestBookingHandler_HandleCancelBookings_UnauthorizedUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		anotherUser = AddUser(tdb.Store, "Unkn", "Unkn", false)
		user        = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel       = AddHotel(tdb.Store, "hotel3", "a", 4, nil)
		room        = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = AddBooking(tdb.Store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	route.Put("/:id", bookingHandler.HandleCancelBookings)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", CreateTokenFromUser(anotherUser))
	resp, err := app.Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	var msg map[string]string
	err = json.NewDecoder(resp.Body).Decode(&msg)
	assert.NoError(t, err)
	assert.Equal(t, msg["error"], "unauthorized")

	// Check if the booking was actually canceled in the database
	updatedBooking, err := tdb.Store.Bookings.GetBookingByID(context.Background(), booking.ID.Hex())
	assert.NoError(t, err)
	assert.False(t, updatedBooking.Canceled)
}
