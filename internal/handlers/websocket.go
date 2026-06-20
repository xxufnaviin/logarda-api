package handlers

import (
	"fmt"
	"logarda/internal/model"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: WebsocketCheckOrigin}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	var request model.WebsocketRequest

	// get username from query params
	request.Username = r.URL.Query().Get("username")
	if request.Username == "" {
		http.Error(w, "Invalid Parameters", http.StatusBadRequest) // return error if not able to parse body
		return
	}

	// add current request as websocket link
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print(err)
		return
	}

	defer conn.Close()

	// append user to list of online users and delete user when connection is lost
	AppendOnlineUser(request.Username)
	defer DeleteOnlineUser(request.Username)

	// goroutine for ws streams
	msgChannel := GetMessageChannel(request.Username)
	go StreamMessages(conn, msgChannel)

	// blocking reader to check websocket connection
	// once this ends, it means websocket connection is lost
	ReadMessages(conn)
}

func WebsocketCheckOrigin(r *http.Request) bool {
	return true // r.Header.Get("Origin") == "https://frontend"
}

func AppendOnlineUser(username string) {
	// create message channel for each user
	msgChannel := make(chan model.Message, 128)

	// append user to list of online users with its own channel
	model.OnlineUsers[username] = msgChannel
}

func DeleteOnlineUser(username string) {
	msgChannel := model.OnlineUsers[username]

	// delete user from online user
	delete(model.OnlineUsers, username)

	// close channel of online user
	close(msgChannel)
}

func GetMessageChannel(username string) chan model.Message {
	return model.OnlineUsers[username]
}

func StreamMessages(conn *websocket.Conn, msgChannel chan model.Message) {
	// for each Message struct in the channel, write to websocket connection
	// waits and blocks when channel is empty
	// only exits when channel is closed or on error
	for msg := range msgChannel {
		err := conn.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}

func ReadMessages(conn *websocket.Conn) {
	// continuously read messages from client
	// will stop if connection is lost
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
