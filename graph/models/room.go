package models

import "time"

type Room struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Type      string    `gorm:"notNull;type:varchar(5);check:type IN ('dm', 'group')" json:"type"` // group or private
	Name      *string   `gorm:"notNull;type:varchar(50)" json:"name"`
	Members   []*User    `gorm:"many2many:roomMembers;" json:"members"`
	CreatedAt time.Time `jsonn:"createdAt"`
}

// type NewRoom struct {
// 	Type      string   `json:"type"`
// 	Name      *string  `json:"name"`
// 	MembersID []string `json:"membersId"`
// }
