package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderChangeLog struct {
	gorm.Model
	ID             uint
	OrderID        uint
	ChangeByUserID uint
	ChangedByUser  Users
	ChangedAt      time.Time
	FieldName      string //status, quantity, item, UnitPrice, Notes, isInternal, vendorID
	oldVal         interface{}
	newVal         interface{}
}
