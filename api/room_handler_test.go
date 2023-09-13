package api

// import (
// 	"encoding/json"
// 	"hotelReservation_golang/types"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/assert"
// )

// func TestRoomHandler_HandleGetRooms(t *testing.T) {
// 	tdb := setup(t)
// 	defer tdb.teardown(t)

// 	app := fiber.New()
// 	roomHandler := NewRoomHandler(tdb.Store)
// 	app.Get("/", roomHandler.HandleGetRooms)

// 	req := httptest.NewRequest("GET", "/rooms", nil)
// 	resp, err := app.Test(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// 	var roomsResp *[]types.Room

// 	err = json.NewDecoder(resp.Body).Decode(&roomsResp)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, roomsResp)
// }
