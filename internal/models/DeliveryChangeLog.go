package models

import (
	"time"
)

type DeliveryChangeLog struct {
	Model
	DeliveryID     uint
	ChangeByUserID uint
	ChangedByUser  Users
	ChangedAt      time.Time
	FieldName      string //scheduledAt, packageType, Notes, Contract, schoolID
	OldVal         interface{}
	NewVal         interface{}
}
