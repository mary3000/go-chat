package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GetUserInfoRequest struct {
	Name string
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	var userReq GetUserInfoRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	Db.Where("username = ?", userReq.Name).First(&user)

	log.Printf("Got user: %v", user)

	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
