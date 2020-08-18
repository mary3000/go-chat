package main

import (
	"bufio"
	"bytes"
	"chat/chat/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const prefix = "-> "

func internalPrefix() string {
	return prefix + fmt.Sprintf("[%v]:[%v] (you) ", currentChatString(), currentNameString())
}

func externalPrefix(user string) string {
	return prefix + fmt.Sprintf("[%v]:[%v] ", currentChatString(), user)
}

func currentNameString() string {
	mu.Lock()
	defer mu.Unlock()

	if currentUser == nil {
		return "_"
	}
	return currentUser.Username
}

func currentChatString() string {
	mu.Lock()
	defer mu.Unlock()

	if currentChat == nil {
		return "_"
	}
	return currentChat.Name
}

var (
	currentUser *server.User
	currentChat *server.Chat
	mu          sync.Mutex

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
	case "r":
		register(parts[1])
	case "l":
		login(parts[1])
	case "cr":
		newChat(parts[1], parts[2:])
	case "cl":
		enterChat(parts[1])
	default:
		userInput <- fmt.Sprint("Fail, unknown command")
	}
}

func enterChat(chat string) {
	mu.Lock()
	userInfo := currentUser
	mu.Unlock()

	if userInfo == nil {
		userInput <- fmt.Sprint("Fail, register/login first")
		return
	}

	//todo: check chat exists
	chatInfo, err := remoteChatInfo(chat)
	if err != nil {
		userInput <- fmt.Sprintf("Fail: %v", err.Error())
		return
	}

	for _, u := range chatInfo.Users {
		userNames[u.ID] = u.Username
	}

	mu.Lock()
	currentChat = chatInfo
	mu.Unlock()

	userInput <- fmt.Sprintf("Chat %v entered", chatInfo.Name)
}

func newChat(chat string, users []string) {
	var userIDs []string
	for _, u := range users {
		userInfo := remoteUserInfo(u)
		userIDs = append(userIDs, strconv.Itoa(int(userInfo.ID)))
	}

	data := server.AddChatRequest{
		Name:  chat,
		Users: userIDs,
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
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		userInput <- fmt.Sprintf("Fail: %v", string(bodyBytes))
	}
}

func login(name string) {
	//todo: check user exists
	userInfo := remoteUserInfo(name)
	if userInfo.ID == 0 {
		userInput <- fmt.Sprint("Fail, user doesn't exist")
	} else {
		mu.Lock()
		currentUser = userInfo
		mu.Unlock()

		userInput <- fmt.Sprint("Logged in")
	}
}

func handleText(text string) {
	//todo: check for chat & name

	mu.Lock()
	chatInfo := currentChat
	userInfo := currentUser
	mu.Unlock()

	if chatInfo == nil || userInfo == nil {
		if chatInfo == nil {
			userInput <- fmt.Sprint("Fail, login into chat first")
		} else {
			userInput <- fmt.Sprint("Fail, login into user first")
		}
		return
	}

	data := server.AddMessageRequest{
		Chat:   strconv.Itoa(int(chatInfo.ID)),
		Author: strconv.Itoa(int(userInfo.ID)),
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
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		userInput <- fmt.Sprintf("Fail: %v", string(bodyBytes))
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

	mu.Lock()
	chatInfo := currentChat
	userInfo := currentUser
	mu.Unlock()

	if chatInfo == nil || userInfo == nil {
		return
	}

	data := server.GetMessagesRequest{
		Chat: strconv.Itoa(int(chatInfo.ID)),
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
				extInput <- fmt.Sprint(externalPrefix(userNames[msgs[i].UserID]) + msgs[i].Text)
			}
			nextMsgIndex = len(msgs)
		}
	} else {
		extInput <- fmt.Sprintf("External fail, status = %v", resp.StatusCode)
	}
}

func remoteChatInfo(name string) (*server.Chat, error) {
	data := server.GetChatInfoRequest{
		Chat: name,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:9000/chatinfo/get", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var chat server.Chat
		err := json.NewDecoder(resp.Body).Decode(&chat)
		if err != nil {
			return nil, err
		}
		return &chat, nil
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return nil, fmt.Errorf(string(bodyBytes))
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
			//fmt.Println("")
			fmt.Println("\r" + text)
		}
	}
}
