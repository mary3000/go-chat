package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

var port = "9000"

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	log.Print("DB opened successfully")

	db.LogMode(true)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Chat{})
	db.AutoMigrate(&Message{})

	mux := http.NewServeMux()

	mux.HandleFunc("/users/add", addUser)
	mux.HandleFunc("/chats/add", addChat)
	mux.HandleFunc("/messages/add", addMessage)

	mux.HandleFunc("/chats/get", getChats)
	mux.HandleFunc("/messages/get", getMessages)

	server := http.Server{Addr: ":" + port, Handler: mux}
	log.Fatal(server.ListenAndServe())
}
