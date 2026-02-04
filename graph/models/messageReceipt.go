package models

import "time"

type MessageReceipt struct {
	MessageID string    `gorm:"primaryKey;type:uuid" json:"messageId"`
	UserID    string    `gorm:"primaryKey;type:uuid" json:"userId"`
	User      *User     `gorm:"foreignKey:userId" json:"user"`
	Status    string    `gorm:"notNull;type:varchar(9);check type IN ('sent', 'delivered', 'read')" json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}
