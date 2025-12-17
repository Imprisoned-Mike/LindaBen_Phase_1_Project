package models

import (
	"time"

	"gorm.io/gorm"
)

type DeliveryChangeLog struct {
	gorm.Model
	ID             uint
	DeliveryID     uint
	ChangeByUserID uint
	ChangedByUser  Users
	ChangedAt      time.Time
	FieldName      string //scheduledAt, packageType, Notes, Contract, schoolID
	OldVal         interface{}
	NewVal         interface{}
}
