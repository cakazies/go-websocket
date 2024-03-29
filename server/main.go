package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-websocket/server/models"
	"github.com/gorilla/websocket"
	"github.com/novalagung/gubrak"
)

type M map[string]interface{}

const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_LEAVE = "Leave"

var (
	conn        *sql.DB
	connections = make([]*WebSocketConnection, 0)
)

type (
	SocketPayload struct {
		Message string
	}

	SocketResponse struct {
		From    string
		Type    string
		Message string
	}

	WebSocketConnection struct {
		*websocket.Conn
		Username string
	}
)

func main() {
	conn := models.Connect()
	if conn == nil {
		log.Println("Tidak Menggunakan Database")
	}
	http.HandleFunc("/", Index)
	http.HandleFunc("/ws", WSwebsocket)
	fmt.Println("Server Starting at :8080")
	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("../client/index.html")
	if err != nil {
		http.Error(w, "Could not open Request file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", content)
}

func WSwebsocket(w http.ResponseWriter, r *http.Request) {
	currentGorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could Not Open Websocket Connection ", http.StatusInternalServerError)
	}

	username := r.URL.Query().Get("username")
	currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username}
	connections = append(connections, &currentConn)
	go handleIO(&currentConn, connections)
}

func handleIO(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error ", fmt.Sprintf("%v", r))
		}
	}()

	broadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				broadcastMessage(currentConn, MESSAGE_LEAVE, "")
				ejectConnection(currentConn)
				return
			}
			log.Println("Error : ", err.Error())
			continue
		}

		broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
	}
}

func ejectConnection(currentConn *WebSocketConnection) {
	filtered, _ := gubrak.Reject(connections, func(each *WebSocketConnection) bool {
		return each == currentConn
	})
	connections = filtered.([]*WebSocketConnection)
}

func broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	for _, eachConn := range connections {
		if eachConn == currentConn {
			continue
		}
		if conn != nil {
			err := models.InsertChat(currentConn.Username, message)
			if err != nil {
				log.Println("Error Insert : ", err)
			}
		}
		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})

	}
}
