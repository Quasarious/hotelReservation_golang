package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testdburi = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type testDB struct {
	db.UserStorage
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	return &testDB{
		UserStorage: db.NewMongoUserStore(client, dbname),
	}
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.UserStorage.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestUserHandler_HandlePostUser_ValidInput(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Post("/", handler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Jameson",
		Email:     "John@gmail.com",
		Password:  "123456!",
	}

	paramsJSON, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(paramsJSON))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody types.User
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.NotNil(t, responseBody.ID)
	assert.Equal(t,
		[]string{params.LastName, params.FirstName, params.Email},
		[]string{responseBody.LastName, responseBody.FirstName, responseBody.Email})

}

func TestUserHandler_HandlePostUser_InvalidInput(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewUserHandler(tdb) // Initialize your UserHandler here
	app.Post("/", handler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "",
		LastName:  "Jameson",
		Email:     "Johngmail.com",
		Password:  "1256!",
	}

	paramsJSON, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(paramsJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var errorResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, errorResponse)
}

func TestUserHandler_HandleGetUser_ValidInput(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := types.User{
		FirstName: "John",
		LastName:  "Jameson",
		Email:     "john@gmail.com",
	}

	handledUser, err := tdb.UserStorage.InsertUser(context.TODO(), &user)
	assert.NoError(t, err)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Get("/:id", handler.HandleGetUser)

	handledUserID := "/" + handledUser.ID.Hex()
	req := httptest.NewRequest("GET", handledUserID, nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody types.User
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, responseBody.FirstName)
	assert.Equal(t, user.LastName, responseBody.LastName)
	assert.Equal(t, user.Email, responseBody.Email)
}

func TestUserHandler_HandleGetUser_UserNotFound(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Get("/:id", handler.HandleGetUser)

	req := httptest.NewRequest("GET", "/non_existent_user", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	// TODO: Change code in HandleGetUser to return http.StatusNotFound
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUserHandler_HandleGetUsers_NoUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Get("/", handler.HandleGetUsers)

	req := httptest.NewRequest("GET", "/", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody []types.User
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Empty(t, responseBody)
}

func TestUserHandler_HandleGetUsers_WithUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	users := []types.User{
		{
			FirstName:         "John",
			LastName:          "Doe",
			Email:             "john@gmail.com",
			EncryptedPassword: "hashed_password_1",
		},
		{
			FirstName:         "Jane",
			LastName:          "Smith",
			Email:             "jane@gmail.com",
			EncryptedPassword: "hashed_password_2",
		},
	}
	for _, user := range users {
		tdb.UserStorage.InsertUser(context.TODO(), &user)
	}

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Get("/", handler.HandleGetUsers)

	req := httptest.NewRequest("GET", "/", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody []types.User
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Len(t, responseBody, len(users))

	for i := range responseBody {
		assert.Equal(t, users[i].FirstName, responseBody[i].FirstName)
		assert.Equal(t, users[i].LastName, responseBody[i].LastName)
		assert.Equal(t, users[i].Email, responseBody[i].Email)
	}
}

func TestUserHandler_HandleDeleteUser_UserExists(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := types.User{
		FirstName: "Mike",
		LastName:  "Jameson",
		Email:     "john@gmail.com",
	}

	handledUser, err := tdb.UserStorage.InsertUser(context.TODO(), &user)
	assert.NoError(t, err)

	userFound, err := tdb.GetUserByID(context.TODO(), handledUser.ID.Hex())
	assert.NoError(t, err)
	assert.NotEmpty(t, userFound)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Delete("/:id", handler.HandleDeleteUser)

	req := httptest.NewRequest("DELETE", "/"+handledUser.ID.Hex(), nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	_, err = tdb.GetUserByID(context.TODO(), handledUser.ID.Hex())
	assert.Error(t, err) // Expecting an error indicating not found
}

func TestUserHandler_HandleDeleteUser_UserNotFound(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Delete("/:id", handler.HandleDeleteUser)

	req := httptest.NewRequest("DELETE", "/non_existent_user", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUserHandler_HandlePutUser_ExistingUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := types.User{
		FirstName: "Mike",
		LastName:  "Jameson",
		Email:     "john@gmail.com",
	}

	handledUser, err := tdb.UserStorage.InsertUser(context.TODO(), &user)
	assert.NoError(t, err)

	userFound, err := tdb.GetUserByID(context.TODO(), handledUser.ID.Hex())
	assert.NoError(t, err)
	assert.NotEmpty(t, userFound)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Put("/:id", handler.HandlePutUser)

	updatedUser := types.UpdateUserParams{
		FirstName: "John",
		LastName:  "Doe",
	}

	updatedUserJSON, err := json.Marshal(updatedUser)
	assert.NoError(t, err)
	req := httptest.NewRequest("PUT", "/"+handledUser.ID.Hex(), bytes.NewBuffer(updatedUserJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"updated: ": handledUser.ID.Hex()}, responseBody)

	userFound, err = tdb.GetUserByID(context.TODO(), handledUser.ID.Hex())
	assert.NoError(t, err)

	assert.Equal(t, updatedUser.FirstName, userFound.FirstName)
	assert.Equal(t, updatedUser.LastName, userFound.LastName)
}

func TestUserHandler_HandlePutUser_InvalidParams(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := types.User{
		FirstName: "Mike",
		LastName:  "Jameson",
		Email:     "john@gmail.com",
	}

	handledUser, err := tdb.UserStorage.InsertUser(context.TODO(), &user)
	assert.NoError(t, err)

	userFound, err := tdb.GetUserByID(context.TODO(), handledUser.ID.Hex())
	assert.NoError(t, err)
	assert.NotEmpty(t, userFound)

	app := fiber.New()
	handler := NewUserHandler(tdb)
	app.Put("/:id", handler.HandlePutUser)

	updatedUser := types.UpdateUserParams{
		FirstName: "J",
		LastName:  "D",
	}

	updatedUserJSON, err := json.Marshal(updatedUser)
	assert.NoError(t, err)
	req := httptest.NewRequest("PUT", "/"+handledUser.ID.Hex(), bytes.NewBuffer(updatedUserJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"updated: ": handledUser.ID.Hex()}, responseBody)

	userFound, err = tdb.GetUserByID(context.TODO(), handledUser.ID.Hex())
	assert.NoError(t, err)

	assert.NotEqual(t, updatedUser.FirstName, userFound.FirstName)
	assert.NotEqual(t, updatedUser.LastName, userFound.LastName)
}
