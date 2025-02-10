package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gihub.com/IsraelTeo/clinic-backend-hackacode-app/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func generateToken(user *model.User) (string, error) {
	expStr := os.Getenv("JWT_EXP")
	iat := time.Now().Unix()
	secret := os.Getenv("API_SECRET")

	if !verifyEnvVariablesToGenerateToken(expStr, secret) {
		log.Println("Environment variables JWT_EXP and API_SECRET are null or invalid")
		return "", fmt.Errorf("environment variables JWT_EXP and API_SECRET are null or invalid")
	}

	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		log.Printf("Error converting JWT_EXP: %v. The default value of 1 hour (%d seconds) will be used", err, exp)
		return "", fmt.Errorf("invalid JWT_EXP value: %v", err)
	}

	claims := Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  iat,
			ExpiresAt: time.Now().Add(time.Second * time.Duration(exp)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error signing the token: %\n", err)
		return "", fmt.Errorf("error signing the token: %v", err)
	}

	return tokenString, nil
}

func verifyEnvVariablesToGenerateToken(exp, secret string) bool {
	return exp != "" || secret != ""
}

func getToken(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	token := parts[1]
	return token, nil
}

func validateToken(c echo.Context) (model.User, error) {
	token, err := getToken(c)
	if err != nil {
		log.Printf("Error retrieving token: %v", err)
		return model.User{}, fmt.Errorf("no token found in request: %w", err)
	}

	// Analiza el token usando jwt.Parse
	jwtToken, err := jwt.Parse(token, validateMethodAndGetSecret)
	if err != nil {
		log.Printf("Token not valid: %v\n", err)
		return model.User{}, fmt.Errorf("invalid token: %w", err)
	}

	// Obtén los claims del token
	userData, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		log.Println("Unable to retrieve payload information or token is invalid")
		return model.User{}, fmt.Errorf("invalid token claims")
	}

	// Verifica que el campo "email" sea un string
	_, ok = userData["email"].(string)
	if !ok {
		log.Println("Email field missing or not a string in token claims")
		return model.User{}, fmt.Errorf("email field is missing or invalid in token claims")
	}

	// Crea el usuario con los datos extraídos del token
	response := model.User{
		Email: userData["email"].(string),
	}

	return response, nil
}

func validateMethodAndGetSecret(token *jwt.Token) (any, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("method not valid")
	}

	// Devuelve el secreto de la variable de entorno
	return []byte(os.Getenv("API_SECRET")), nil
}
