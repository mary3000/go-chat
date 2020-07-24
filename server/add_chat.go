package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

type addChatRequest struct {
	Name  string
	Users []string
}

func addChat(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	var chatReq addChatRequest
	err := json.NewDecoder(r.Body).Decode(&chatReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat := &Chat{Name: chatReq.Name}

	users := []User{}
	for _, u := range chatReq.Users {
		id, err := strconv.Atoi(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		users = append(users, User{
			Model: gorm.Model{ID: uint(id)},
		})
	}
	chat.Users = users

	res := db.Create(chat)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Added chat: %v", chatReq)

	//todo: check for user existence

	_, _ = w.Write([]byte(fmt.Sprintf("Chat id: %v", chat.ID)))
}
