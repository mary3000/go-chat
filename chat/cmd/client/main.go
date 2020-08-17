package main

import (
	"bufio"
	"bytes"
	"chat/chat/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const prefix = "-> "

func internalPrefix() string {
	return prefix + fmt.Sprintf("[%v]:[%v] ", currentChatString(), currentNameString())
}

func externalPrefix(user string) string {
	return prefix + fmt.Sprintf("[%v]:[%v] ", currentChatString(), user)
}

func currentNameString() string {
	if currentUser == nil {
		return ""
	}
	return currentUser.Username
}

func currentChatString() string {
	if currentChat == nil {
		return ""
	}
	return currentChat.Name
}

var (
	currentUser  *server.User
	currentChat  *server.Chat
	nextMsgIndex = 0
	userNames    = make(map[uint]string)
)

var userInput = make(chan string)
var extInput = make(chan string)

func register(name string) {
	data := server.AddUserRequest{
		Username: name,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/users/add", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		userInput <- fmt.Sprintf("User %v registered\n", name)
	} else {
		userInput <- fmt.Sprintf("Fail, status = %v", resp.StatusCode)
	}
}

func isCmd(text string) bool {
	return strings.HasPrefix(text, "/")
}

func handleCmd(text string) {
	parts := strings.Split(text, " ")
	cmd := parts[0]

	switch cmd {
	case "register":
		register(parts[1])
	case "login":
		login(parts[1])
	case "newchat":
		newChat(parts[1], parts[2:])
	case "enterchat":
		enterChat(parts[1])
	}
}

func enterChat(chat string) {
	if currentUser == nil {
		userInput <- fmt.Sprint("Fail, register/login first")
		return
	}

	//todo: check chat exists
	currentChat = remoteChatInfo(chat)

	for _, u := range currentChat.Users {
		userNames[u.ID] = u.Username
	}

	userInput <- fmt.Sprint("Chat entered")
}

func newChat(chat string, users []string) {
	data := server.AddChatRequest{
		Name:  chat,
		Users: users,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/chats/add", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		userInput <- fmt.Sprintf("Chat %v registered", chat)
	} else {
		userInput <- fmt.Sprintf("Fail, status = %v", resp.StatusCode)
	}
}

func login(name string) {
	//todo: check user exists
	currentUser = remoteUserInfo(name)
	userInput <- fmt.Sprint("Logged in")
}

func handleText(text string) {
	//todo: check for chat & name

	if currentChat == nil || currentUser == nil {
		userInput <- fmt.Sprint("Fail, register/login first")
		return
	}

	data := server.AddMessageRequest{
		Chat:   strconv.Itoa(int(currentChat.ID)),
		Author: strconv.Itoa(int(currentUser.ID)),
		Text:   text,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/messages/add", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		//userInput <- fmt.Sprint("")
	} else {
		userInput <- fmt.Sprintf("Fail, status = %v", resp.StatusCode)
	}
}

func handleUserInput() {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if isCmd(text) {
			handleCmd(text[1:])
		} else {
			handleText(text)
		}

	}
}

func getMessages() {
	//todo: check for chat

	if currentChat == nil || currentUser == nil {
		return
	}

	data := server.GetMessagesRequest{
		Chat: strconv.Itoa(int(currentChat.ID)),
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/messages/get", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var msgs []server.Message
		err := json.NewDecoder(resp.Body).Decode(&msgs)

		if err != nil {
			panic(err)
		}

		if len(msgs) > nextMsgIndex {
			for i := nextMsgIndex; i < len(msgs); i++ {
				//if msgs[i].UserID != currentUser.ID {
				extInput <- fmt.Sprint(externalPrefix(userNames[msgs[i].UserID]) + msgs[i].Text)
				//}
			}
			nextMsgIndex = len(msgs)
		}
	} else {
		extInput <- fmt.Sprintf("External fail, status = %v", resp.StatusCode)
	}
}

func remoteChatInfo(name string) *server.Chat {
	data := server.GetChatInfoRequest{
		Chat: name,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/chatinfo/get", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var chat server.Chat
		err := json.NewDecoder(resp.Body).Decode(&chat)
		if err != nil {
			panic(err)
		}
		return &chat
	} else {
		panic(resp.StatusCode)
	}
}

func remoteUserInfo(name string) *server.User {
	data := server.GetUserInfoRequest{
		Name: name,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/userinfo/get", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var user server.User
		err := json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			panic(err)
		}
		return &user
	} else {
		panic(resp.StatusCode)
	}
}

func pollChat() {
	for {
		getMessages()
		time.Sleep(time.Millisecond * 300)
	}
}

func main() {
	fmt.Println("Go Chat")
	fmt.Println("---------------------")

	go handleUserInput()
	go pollChat()

	for {
		fmt.Print(internalPrefix())
		select {
		case text := <-userInput:
			fmt.Println(text)
		case text := <-extInput:
			fmt.Println("")
			fmt.Println(text)
		}
	}
}
