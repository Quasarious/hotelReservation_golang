package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
	"os"
	"time"
)

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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid email or password")
		}
		return err
	}

	if !types.IsPasswordValid(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("invalid email or password")
	}
	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(4 * time.Hour)
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   validTill.Unix(),
	}
	claims["id"] = user.ID
	claims["email"] = user.Email

	secret := os.Getenv("JWT_SECRET")
	strToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		fmt.Println("error creating token")
	}
	return strToken
}
