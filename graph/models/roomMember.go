package models

import "time"

type RoomMember struct {
	RoomID            string    `gorm:"primaryKey;type:uuid" json:"roomId"`
	UserID            string    `gorm:"primaryKey;type:uuid" json:"userId"`
	LastReadMessageID string    `gorm:"type:uuid" json:"lastReadMessageId"`
	LastReadAt        time.Time `json:"lastReadAt"`
}
