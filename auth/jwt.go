package auth

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey = []byte("c2VjcmV0a2V5c3VwZXItaGFzaA==")
)

// GenerateToken creates a new JWT token
func GenerateToken(userID int) (string, error) {
	log.Printf("Generating token for userID: %d", userID)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	log.Printf("Token generated successfully: %s", tokenString)
	return tokenString, nil
}

// ParseToken parses a JWT token and returns the claims
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	log.Printf("Received token in Parse Token function: %s", tokenString)

	// Check if token starts with "Bearer " and remove it
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("Token parsed successfully. Claims: %v", claims)
		return claims, nil
	}

	log.Println("Invalid token")
	return nil, errors.New("invalid token")
}
