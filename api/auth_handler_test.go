package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func makeTestUser(t *testing.T, tdb *db.Store) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "John",
		LastName:  "Jameson",
		Email:     "John@gmail.com",
		Password:  "123456!",
	})

	if err != nil {
		log.Fatal(err)
	}

	completeUser, err := tdb.Users.InsertUser(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}

	return completeUser
}

func TestAuthHandler_HandleAuthenticate_Success(t *testing.T) {
	tdb := setup(t)
	insertedUser := makeTestUser(t, tdb.Store)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewAuthHandler(tdb.Store)
	app.Post("/auth", handler.HandleAuthenticate)

	b, err := json.Marshal(AuthParams{
		Email:    "John@gmail.com",
		Password: "123456!",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200, but got %d", err)
	}

	var responseBody AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "John@gmail.com", responseBody.User.Email)
	assert.NotEmpty(t, responseBody.Token)
	reflect.DeepEqual(insertedUser, responseBody.User)
}

func TestAuthHandler_HandleAuthenticate_Failed(t *testing.T) {
	tdb := setup(t)
	makeTestUser(t, tdb.Store)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewAuthHandler(tdb.Store)
	app.Post("/auth", handler.HandleAuthenticate)

	b, err := json.Marshal(AuthParams{
		Email:    "John@gmail.com",
		Password: "16!",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected http status of 401, but got: %v", err)
	}

	var responseBody AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.Error(t, err)
	assert.Empty(t, responseBody)
	assert.Empty(t, responseBody.Token)
}

func TestAuthHandler_HandleAuthenticate_InvalidJSON(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewAuthHandler(tdb.Store)
	app.Post("/auth", handler.HandleAuthenticate)

	// Invalid JSON request body
	b := []byte(`invalid-json`)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var responseBody AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.Error(t, err)
	assert.Empty(t, responseBody)
	assert.Empty(t, responseBody.Token)
}

func TestAuthHandler_HandleAuthenticate_MissingFields(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewAuthHandler(tdb.Store)
	app.Post("/auth", handler.HandleAuthenticate)

	// Missing email
	b, err := json.Marshal(AuthParams{
		Password: "123456!",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	// TODO : fix status code
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Missing password
	b, err = json.Marshal(AuthParams{
		Email: "John@gmail.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseBody AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.Error(t, err)
	assert.Empty(t, responseBody)
	assert.Empty(t, responseBody.Token)
}
