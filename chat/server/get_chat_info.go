package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GetChatInfoRequest struct {
	Chat string
}

func GetChatInfo(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	var chatsReq GetChatInfoRequest
	err := json.NewDecoder(r.Body).Decode(&chatsReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var chat Chat
	Db.Model(&chat).Association("Users")
	Db.Where("name = ?", chatsReq.Chat).First(&chat)
	Db.Model(&chat).Association("Users").Find(&chat.Users)

	log.Printf("Got chat: %v ", chat)

	js, err := json.Marshal(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
