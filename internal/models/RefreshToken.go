package models

import (
	"time"

	"LindaBen_Phase_1_Project/internal/db"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func CreateRefreshToken(userID uint) (*RefreshToken, error) {
	rt := RefreshToken{
		UserID:    userID,
		Token:     uuid.NewString(), // random string
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	return &rt, db.Db.Create(&rt).Error
}

func ValidateRefreshToken(token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := db.Db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&rt).Error
	return &rt, err
}

func DeleteRefreshToken(token string) error {
	result := db.Db.Where("token = ?", token).Delete(&RefreshToken{})

	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return result.Error
}
