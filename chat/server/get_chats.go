package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type getChatsRequest struct {
	User string
}

type chatDescriptor struct {
	Name      string
	ID        uint
	CreatedAt time.Time
}

func GetChats(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	var chatsReq getChatsRequest
	err := json.NewDecoder(r.Body).Decode(&chatsReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(chatsReq.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := Db.Table("chats").
		Select("chats.name, chats.id, chats.created_at, max(messages.created_at) as update_time").
		Joins("join chat_users on chat_users.chat_id = chats.id").
		Where("chat_users.user_id = ?", userID).
		Joins("left join messages on messages.chat_id = chats.id").
		Group("chats.id").
		Order("update_time DESC").
		Rows()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	chats := []chatDescriptor{}
	for rows.Next() {
		var chatDesc chatDescriptor
		var updateTime string
		_ = rows.Scan(&chatDesc.Name, &chatDesc.ID, &chatDesc.CreatedAt, &updateTime)
		chats = append(chats, chatDesc)
		log.Printf("row: %v", chatDesc)
	}

	log.Printf("Found chats: %v", chats)

	js, err := json.Marshal(chats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
