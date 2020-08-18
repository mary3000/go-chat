package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var testDbName = "test_db"

func Server() {
	dbName = testDbName
	_ = os.Remove(testDbName)
	main()
}

// Methods below are based on curl-to-Go: https://mholt.github.io/curl-to-go

func testUsersAdd(name string) error {
	type Payload struct {
		Username string `json:"username"`
	}

	data := Payload{
		Username: name,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/users/add", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func testChatsAdd(name string, users []string) error {
	type Payload struct {
		Name  string   `json:"name"`
		Users []string `json:"users"`
	}

	data := Payload{
		Name:  name,
		Users: users,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/chats/add", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func testMessagesAdd(chat string, name string, text string) error {
	type Payload struct {
		Chat   string `json:"chat"`
		Author string `json:"author"`
		Text   string `json:"text"`
	}

	data := Payload{
		Chat:   chat,
		Author: name,
		Text:   text,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/messages/add", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func testChatsGet(name string) error {
	type Payload struct {
		User string `json:"user"`
	}

	data := Payload{
		User: name,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/chats/get", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func testMessagesGet(chat string) error {
	type Payload struct {
		Chat string `json:"chat"`
	}

	data := Payload{
		Chat: chat,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/messages/get", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type Request int

const (
	usersAdd Request = iota
	chatsAdd
	messagesAdd
	chatsGet
	messagesGet
)

func TestAPI(t *testing.T) {
	go Server()

	time.Sleep(time.Millisecond * 300) // wait for db startup

	requests := []struct {
		reqType Request
		name    string
		chat    string
		text    string
		users   []string
	}{
		{
			reqType: usersAdd,
			name:    "user1",
		},
		{
			reqType: usersAdd,
			name:    "user2",
		},
		{
			reqType: chatsAdd,
			chat:    "chat1",
			users:   []string{"1", "2"},
		},
		{
			reqType: messagesAdd,
			chat:    "chat1",
			name:    "user1",
			text:    "hello",
		},
		{
			reqType: chatsGet,
			name:    "1",
		},
		{
			reqType: messagesGet,
			chat:    "1",
		},
	}

	for _, r := range requests {
		var err error
		switch r.reqType {
		case usersAdd:
			err = testUsersAdd(r.name)
		case chatsAdd:
			err = testChatsAdd(r.chat, r.users)
		case messagesAdd:
			err = testMessagesAdd(r.chat, r.name, r.text)
		case chatsGet:
			err = testChatsGet(r.name)
		case messagesGet:
			err = testMessagesGet(r.chat)
		default:
			log.Fatal("unknown request type")
		}
		if err != nil {
			t.Errorf("err = %v", err)
		}
	}
}
