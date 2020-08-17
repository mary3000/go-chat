package server

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// id, created_at - defined in gorm.Model
type User struct {
	gorm.Model

	Username string `gorm:"unique"`
}

// id, created_at - defined in gorm.Model
type Chat struct {
	gorm.Model

	Name  string `gorm:"unique"`
	Users []User `gorm:"many2many:chat_users"`
}

// id, created_at - defined in gorm.Model
type Message struct {
	gorm.Model

	ChatID uint
	UserID uint // author
	Text   string
}

var Db *gorm.DB
