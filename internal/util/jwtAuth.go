package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthAdmin checks for a valid token with the "admin" role.
func JWTAuthAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		err = ValidateAdminRoleJWT(context)
		if err != nil {
			context.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admin role required"})
			context.Abort()
			return
		}
		context.Next()
	}
}

// JWTAuthSchool checks for a valid token with "school_admin" or "admin" role.
func JWTAuthSchool() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		err = ValidateSchoolAdminRoleJWT(context)
		if err != nil {
			context.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: School admin or admin role required"})
			context.Abort()
			return
		}
		context.Next()
	}
}

// JWTAuthVendor checks for a valid token with "vendor_admin" or "admin" role.
func JWTAuthVendor() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		err = ValidateVendorAdminRoleJWT(context)
		if err != nil {
			context.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Vendor admin or admin role required"})
			context.Abort()
			return
		}
		context.Next()
	}
}
