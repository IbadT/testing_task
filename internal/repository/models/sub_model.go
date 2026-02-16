package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ServiceName string     `gorm:"type:varchar(50);not null"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`
	Price       int        `gorm:"type:int;not null"`
	StartDate   time.Time  `gorm:"type:date;not null"`
	EndDate     *time.Time `gorm:"type:date;null"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
