package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
	"os"
	"time"
)

const tokenExpirationTime = 4 * time.Hour

type AuthHandler struct {
	store db.Store
}

func NewAuthHandler(storage *db.Store) *AuthHandler {
	return &AuthHandler{
		store: *storage,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.store.Users.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email or password!")
	}

	if !types.IsPasswordValid(user.EncryptedPassword, params.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password!")
	}

	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(tokenExpirationTime)
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   validTill.Unix(),
	}

	secret := os.Getenv("JWT_SECRET")
	strToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		fmt.Println("error creating token")
	}
	return strToken
}
