package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("JWT Authentication")

	token, ok := c.GetReqHeaders()["X-Access-Token"]
	if !ok {
		return fmt.Errorf("X-Access-Token not found")
	}
	claims, err := validateToken(token)
	if err != nil {
		return err
	}
	expiresFloat := claims["exp"].(float64)

	if int64(expiresFloat)-time.Now().Unix() < 0 {
		return fmt.Errorf("token expired")
	}

	// check token expiration
	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	fmt.Println("paresJWTToken")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("unexpected signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("unauthorized")
}
