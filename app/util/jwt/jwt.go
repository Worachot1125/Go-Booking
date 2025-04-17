package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const (
	ContextUser = "AUTH_USER"
)

type UserPayload struct {
	UserID uuid.UUID `json:"user_id"`
}
type SigningMethodHMAC = jwt.SigningMethodHMAC

// VerifyToken verifies the JWT token and returns the claims as a map or an error
func VerifyToken(raw string) (map[string]any, error) {

	token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method conforms to expected HMAC method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

		// Use the secret key from the configuration
		secret := []byte(viper.GetString("TOKEN_SECRET_USER"))
		return secret, nil
	})

	if err != nil {
		// Return a detailed error if token parsing fails
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	// Return a generic invalid token error if the token is not valid
	return nil, errors.New("invalid token")
}

func CreateToken(claims jwt.MapClaims, secretKey string) (string, error) {
	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	secret := []byte(secretKey)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
