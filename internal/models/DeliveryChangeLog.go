package models

import (
	"time"
)

type DeliveryChangeLog struct {
	Model
	DeliveryID     uint      `json:"deliveryId"`
	ChangeByUserID uint      `json:"changeByUserId"`
	ChangedByUser  Users     `gorm:"foreignKey:ChangeByUserID" json:"changedByUser"`
	ChangedAt      time.Time `json:"changedAt"`
	FieldName      string    `json:"fieldName"` // scheduledAt, packageType, Notes, Contract, schoolID
	OldValue       string    `json:"oldValue"`
	NewValue       string    `json:"newValue"`
}
