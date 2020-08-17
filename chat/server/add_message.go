package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type AddMessageRequest struct {
	Chat   string
	Author string
	Text   string
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	var msgReq AddMessageRequest
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(msgReq.Chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	authorID, err := strconv.Atoi(msgReq.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := &Message{
		ChatID: uint(chatID),
		UserID: uint(authorID),
		Text:   msgReq.Text}

	res := Db.Create(msg)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Added message: %v", msgReq)

	//todo: check for user and chat existence

	_, _ = w.Write([]byte(fmt.Sprintf("Message id: %v", msg.ID)))
}
