package models

import (
	"time"
)

type OrderChangeLog struct {
	Model
	OrderID        uint      `json:"orderId"`
	ChangeByUserID uint      `json:"changeByUserId"`
	ChangedByUser  User      `gorm:"foreignKey:ChangeByUserID" json:"changedByUser"`
	ChangedAt      time.Time `json:"changedAt"`
	FieldName      string    `json:"fieldName"` // status, quantity, item, UnitPrice, Notes, isInternal, vendorID
	OldValue       string    `json:"oldValue"`
	NewValue       string    `json:"newValue"`
}
