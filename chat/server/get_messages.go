package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type GetMessagesRequest struct {
	Chat string
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	var msgReq GetMessagesRequest
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

	var msgs []Message
	Db.Where("chat_id = ?", chatID).Find(&msgs).Order("created_at ASC")

	log.Printf("Found messages for %v", msgReq)

	js, err := json.Marshal(msgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
