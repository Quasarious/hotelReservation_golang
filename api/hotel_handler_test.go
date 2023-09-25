package api

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservation_golang/middleware"
	"hotelReservation_golang/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHotelHandler_HandleGetHotels_Correct(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var hotels []*types.Hotel
	hotels = append(hotels, AddHotel(tdb.Store, "MoscowHotel", "Moscow", 3.4, []primitive.ObjectID{}))
	hotels = append(hotels, AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{}))
	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		room         = AddRoom(tdb.Store, "Large", false, 99.9, hotels[0].ID)
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	hotels[0].Rooms = append(hotels[0].Rooms, room.ID)
	route.Get("/", hotelHandler.HandleGetHotels)
	req := httptest.NewRequest(http.MethodGet, "/hotels", nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var queryResp *HotelsSourceResp
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	fmt.Printf("%+v\n", queryResp)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(hotels[0], queryResp.Hotels[0]))
}

func TestHotelHandler_HandleGetHotel_CorrectData(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel        = AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{})
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	route.Get("/", hotelHandler.HandleGetHotel)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", hotel.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var queryResp *types.Hotel
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(hotel, queryResp))
}

func TestHotelHandler_HandleGetHotel_InvalidID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		_            = AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{})
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	route.Get("/", hotelHandler.HandleGetHotel)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", "213sda"), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var queryResp *types.Hotel
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	assert.NoError(t, err)
	assert.Empty(t, queryResp)
}

func TestHotelHandler_HandleGetHotel_NoHotel(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		_            = AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{})
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	route.Get("/", hotelHandler.HandleGetHotel)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", user.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var hotelResp *types.Hotel
	err = json.NewDecoder(resp.Body).Decode(&hotelResp)
	assert.NoError(t, err)
	assert.Empty(t, hotelResp)
}

func TestHotelHandler_HandleGetRooms_CorrectData(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel        = AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{})
		room         = AddRoom(tdb.Store, "Large", false, 99.9, hotel.ID)
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	hotel.Rooms = append(hotel.Rooms, room.ID)
	route.Get("/:id", hotelHandler.HandleGetRooms)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/rooms", hotel.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var roomResp []*types.Room
	err = json.NewDecoder(resp.Body).Decode(&roomResp)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(roomResp[0], room))
}

func TestHotelHandler_HandleGetRooms_InvalidPath(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel        = AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{})
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	route.Get("/", hotelHandler.HandleGetRooms)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/rooms", hotel.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var roomResp []*types.Room
	err = json.NewDecoder(resp.Body).Decode(&roomResp)
	assert.Error(t, err)
	assert.Empty(t, roomResp)
}

func TestHotelHandler_HandleDeleteHotel_Success(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel        = AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{})
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	route.Delete("/", hotelHandler.HandleDeleteHotel)
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", hotel.ID.Hex()), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var roomResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&roomResp)
	assert.NoError(t, err)
	assert.Equal(t, roomResp["deleted:"], "Hotel with ID: "+hotel.ID.Hex())
}

func TestHotelHandler_HandleDeleteHotel_InvalidID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user         = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		_            = AddHotel(tdb.Store, "KazanHotel", "Kazan", 4.4, []primitive.ObjectID{})
		app          = fiber.New()
		route        = app.Group("/:id", middleware.JWTAuthentication(tdb.Store.Users))
		hotelHandler = NewHotelHandler(tdb.Store)
	)
	route.Delete("/", hotelHandler.HandleDeleteHotel)
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", "sdgsd"), nil)
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var roomResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&roomResp)
	assert.NoError(t, err)
	assert.Equal(t, roomResp["error"], "No hotels found with such ID")
}
