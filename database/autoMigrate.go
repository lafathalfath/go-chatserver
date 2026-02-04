package database

import (
	"github.com/lafathalfath/go-chatserver/graph/models"

	"gorm.io/gorm"
)

func autoMigrate (db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Message{},
		&models.RoomMember{},
		&models.MessageReceipt{},
	)
	if err != nil {
		panic(err)
	}
}