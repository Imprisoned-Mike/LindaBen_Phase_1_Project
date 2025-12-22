package util

import (
	"LindaBen_Phase_1_Project/internal/models"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

// retrieve JWT key from .env file
var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// generate JWT token
func GenerateAccessToken(user models.Users) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID, // subject (standard claim)
		"role": user.Roles,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 day token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// validate Admin role
func ValidateAdminRoleJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token provided")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return errors.New("role claim is missing or not a string")
	}

	if role == "admin" {
		return nil
	}
	return errors.New("user is not an admin")
}

// validate School Admin role (or admin)
func ValidateSchoolAdminRoleJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token provided")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return errors.New("role claim is missing or not a string")
	}

	if role == "school_admin" || role == "admin" {
		return nil
	}
	return errors.New("user is not a school admin or admin")
}

// validate Vendor Admin role (or admin)
func ValidateVendorAdminRoleJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token provided")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return errors.New("role claim is missing or not a string")
	}

	if role == "vendor_admin" || role == "admin" {
		return nil
	}
	return errors.New("user is not a vendor admin or admin")
}

// fetch user details from the token
func CurrentUser(context *gin.Context) models.Users {
	err := ValidateJWT(context)
	if err != nil {
		return models.Users{}
	}
	token, _ := getToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	// Use "sub" (subject) claim for user ID, which is standard
	userIdFloat, ok := claims["sub"].(float64)
	if !ok {
		return models.Users{}
	}
	userId := uint(userIdFloat)

	var user models.Users
	// models.GetUser expects an int, so we cast
	err = models.GetUser(&user, int(userId))
	if err != nil {
		return models.Users{}
	}
	return user
}

// check token validity
func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

// extract token from request Authorization header
func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
