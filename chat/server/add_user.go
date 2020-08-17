package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AddUserRequest struct {
	Username string
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		http.Error(w, fmt.Sprintf("Content-type: expected %v, got %v", expectedContentType, contentType),
			http.StatusBadRequest)
		return
	}

	expectedMethod := "POST"
	if r.Method != expectedMethod {
		http.Error(w, fmt.Sprintf("Method: expected %v, got %v", expectedMethod, r.Method), http.StatusBadRequest)
	}

	var u AddUserRequest
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &User{Username: u.Username}
	res := Db.Create(user)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Added user: %v", u)

	_, _ = w.Write([]byte(fmt.Sprintf("User id: %v", user.ID)))
}
