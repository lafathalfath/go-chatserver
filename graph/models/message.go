package models

import (
	"time"
)

type Message struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Content   string    `gorm:"notNull;type:text" json:"content"`
	SenderID  string    `gorm:"notNull;type:uuid" json:"senderId"`
	Sender    *User      `gorm:"foreignKey:senderId" json:"sender"`
	RoomID    string    `gorm:"notNull;type:uuid" json:"roomId"`
	CreatedAt time.Time `json:"created_at"`
}

// type NewMessage struct {
// 	Content  string `json:"content"`
// 	RoomID string `json:"roomId"`
// }
