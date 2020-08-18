package main

import (
	"chat/chat/server"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

var port = "9000"

func main() {
	var err error
	server.Db, err = gorm.Open("sqlite3", "db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer server.Db.Close()

	log.Print("chat.Db opened successfully")

	server.Db.LogMode(true)
	server.Db.AutoMigrate(&server.User{})
	server.Db.AutoMigrate(&server.Chat{})
	server.Db.AutoMigrate(&server.Message{})

	mux := http.NewServeMux()

	mux.HandleFunc("/users/add", server.AddUser)
	mux.HandleFunc("/chats/add", server.AddChat)
	mux.HandleFunc("/messages/add", server.AddMessage)

	mux.HandleFunc("/chats/get", server.GetChats)
	mux.HandleFunc("/messages/get", server.GetMessages)

	// for cli
	mux.HandleFunc("/chatinfo/get", server.GetChatInfo)
	mux.HandleFunc("/userinfo/get", server.GetUserInfo)

	chat := http.Server{Addr: ":" + port, Handler: mux}
	log.Fatal(chat.ListenAndServe())
}
