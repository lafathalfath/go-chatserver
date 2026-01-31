package models

type User struct {
	ID       string `gorm:"type:uuid;primaryKey" json:"id"`
	Name     string `gorm:"notNull;type:varchar(50)" json:"name"`
	Email    string `gorm:"notNull;unique;type:varchar(100)" json:"email"`
	Password string `gorm:"notNull;type:varchar(255)" json:"password"`
	Rooms    []*Room `gorm:"many2many:roomMembers" json:"rooms"`
}

// type NewUser struct {
// 	Name                 string `json:"name"`
// 	Email                string `json:"email"`
// 	Password             string `json:"password"`
// 	PasswordConfirmation string `json:"passwordConfirmation"`
// }

// type Login struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }
