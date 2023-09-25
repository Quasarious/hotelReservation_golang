package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservation_golang/middleware"
	"hotelReservation_golang/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRoomHandler_HandleGetRooms(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	hotel := AddHotel(tdb.Store, "hotel3", "a", 4, []primitive.ObjectID{})
	room := AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)
	hotel.Rooms = append(hotel.Rooms, room.ID)
	app := fiber.New()
	roomHandler := NewRoomHandler(tdb.Store)
	app.Get("/", roomHandler.HandleGetRooms)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var roomsResp *RoomsSourceResp

	err = json.NewDecoder(resp.Body).Decode(&roomsResp)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(room, roomsResp.Rooms[0]))
}

func TestRoomHandler_HandleBookRoom_ValidData(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		user  = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel = AddHotel(tdb.Store, "hotel3", "a", 4, []primitive.ObjectID{})
		room  = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from     = time.Now().AddDate(0, 0, 1)
		till     = time.Now().AddDate(0, 0, 5)
		bookRoom = BookRoomParams{
			NumPersons: 2,
			FromDate:   from,
			TillDate:   till,
		}
		app         = fiber.New()
		route       = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		roomHandler = NewRoomHandler(tdb.Store)
	)
	hotel.Rooms = append(hotel.Rooms, room.ID)
	b, err := json.Marshal(&bookRoom)
	assert.NoError(t, err)

	route.Post("/book/:id", roomHandler.HandleBookRoom)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/book/%s", room.ID.Hex()), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var queryResp *types.Booking
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	assert.NoError(t, err)
	assert.Equal(t, bookRoom.FromDate, queryResp.FromDate)
	assert.Equal(t, bookRoom.TillDate, queryResp.TillDate)
	assert.Equal(t, bookRoom.NumPersons, queryResp.NumPersons)
	assert.Equal(t, room.ID, queryResp.RoomID)
	assert.Equal(t, user.ID, queryResp.UserID)
}

func TestRoomHandler_HandleBookRoom_InvalidRoomID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		user  = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel = AddHotel(tdb.Store, "hotel3", "a", 4, []primitive.ObjectID{})
		room  = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from     = time.Now().AddDate(0, 0, 1)
		till     = time.Now().AddDate(0, 0, 5)
		bookRoom = BookRoomParams{
			NumPersons: 2,
			FromDate:   from,
			TillDate:   till,
		}
		app         = fiber.New()
		route       = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		roomHandler = NewRoomHandler(tdb.Store)
	)
	hotel.Rooms = append(hotel.Rooms, room.ID)
	b, err := json.Marshal(&bookRoom)
	assert.NoError(t, err)

	route.Post("/book/:id", roomHandler.HandleBookRoom)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/book/%s", "asdasd"), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var queryResp *types.Booking
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	assert.NoError(t, err)
	assert.Empty(t, queryResp)
}

func TestRoomHandler_HandleBookRoom_InvalidDateTime(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		user  = AddUser(tdb.Store, "Iak", "Vapvapvv", false)
		hotel = AddHotel(tdb.Store, "hotel3", "a", 4, []primitive.ObjectID{})
		room  = AddRoom(tdb.Store, "small", true, 4.99, hotel.ID)

		from     = time.Now().AddDate(0, 0, -1)
		till     = time.Now().AddDate(0, 0, 5)
		bookRoom = BookRoomParams{
			NumPersons: 2,
			FromDate:   from,
			TillDate:   till,
		}
		app         = fiber.New()
		route       = app.Group("/", middleware.JWTAuthentication(tdb.Store.Users))
		roomHandler = NewRoomHandler(tdb.Store)
	)
	hotel.Rooms = append(hotel.Rooms, room.ID)
	b, err := json.Marshal(&bookRoom)
	assert.NoError(t, err)

	route.Post("/book/:id", roomHandler.HandleBookRoom)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/book/%s", "asdasd"), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var queryResp *types.Booking
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	assert.NoError(t, err)
	assert.Empty(t, queryResp)
}
