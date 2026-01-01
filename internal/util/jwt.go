package util

import (
	"LindaBen_Phase_1_Project/internal/models"
	"errors"
	"fmt"
	"net/http"
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
		"sub":   user.ID, // subject (standard claim)
		"roles": user.Roles,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 day token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(privateKey)
}

// JWTAuth checks for a valid token with the "admin" role.
func JWTAuth(allowedRoles ...string) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		err = ValidateRoleJWT(context, allowedRoles...)
		if err != nil {
			context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
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

// validate if user has one of the allowed roles
func ValidateRoleJWT(context *gin.Context, allowedRoles ...string) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token provided")
	}

	roles, ok := claims["roles"].(string)
	if !ok {
		return errors.New("role claim is missing or not a string")
	}

	rolesArray := models.ParseRoles(roles)
	fmt.Println("User roles from token:", rolesArray)

	for _, allowedRole := range allowedRoles {
		if strings.Contains(roles, allowedRole) {
			return nil
		}
	}

	return fmt.Errorf("access denied: requires one of the roles %v", allowedRoles)
}

// fetch user details from the token
func CurrentUser(context *gin.Context) *models.Users {
	err := ValidateJWT(context)
	if err != nil {
		return nil
	}
	token, _ := getToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	// Use "sub" (subject) claim for user ID, which is standard
	userIdFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil
	}
	userId := uint(userIdFloat)

	var user models.Users
	// models.GetUser expects an int, so we cast
	err = models.GetUser(&user, int(userId))
	if err != nil {
		return nil
	}
	return &user
}

// check token validity
func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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
