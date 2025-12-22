package models

import (
	"time"
)

type OrderChangeLog struct {
	Model
	OrderID        uint
	ChangeByUserID uint
	ChangedByUser  Users
	ChangedAt      time.Time
	FieldName      string //status, quantity, item, UnitPrice, Notes, isInternal, vendorID
	OldVal         interface{}
	NewVal         interface{}
}
