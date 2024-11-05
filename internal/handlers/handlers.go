package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

// channels
var websocketChannel = make(chan WebsocketPayload)
var clients = make(map[WebSocketConnection]string)

// views
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

// websockets upgrader
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Home
func Home(w http.ResponseWriter, r *http.Request) {
	copyFileToFile()
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// websocket sections
type WebSocketConnection struct {
	*websocket.Conn
}
type WebsocketJsonResponse struct {
	Action       string `json:"action"`
	Message      string `json:"message"`
	Message_Type string `json:"message_type"`
}
type WebsocketPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"_"`
}

func WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebsocketEndpoint 1: ", err)
	}
	log.Println("WebsocketEndpoint: Client is connected to endpoint !")

	var response WebsocketJsonResponse
	// response.Message = `<em><small>Connected to the server</small></em>`
	response.Message = `Connected to the server`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("WebsocketEndpoint 2: ", err)
	}

	go ListenForWebsocket(&conn)
}

func ListenForWebsocket(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error, ListenForWebsocket 1: ", fmt.Sprintf("%v", r))
		}
	}()

	var payload WebsocketPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Println("Error, ListenForWebsocket 2: ", err)
		} else {
			payload.Conn = *conn
			websocketChannel <- payload
			log.Println("websocketChannel")
		}
	}
}

func ListenToWebsocketChannel() {
	var response WebsocketJsonResponse

	for {
		messageFromWSChannel := <-websocketChannel

		response.Action = "Message Received !"
		response.Message = fmt.Sprintf("Some message, and action was %s", messageFromWSChannel.Action)
		boardcastToAllUser(response)
	}
}

func boardcastToAllUser(response WebsocketJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Error, boardcastToAllUser: ", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

// renderPages
func renderPage(w http.ResponseWriter, givenTemplate string, data jet.VarMap) error {
	view, err := views.GetTemplate(givenTemplate)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func copyFileToFile() {
	// Open file1 for reading
	srcFile, err := os.Open("./html/index.html")
	if err != nil {
		log.Fatalf("Failed to open index.html: %v", err)
	}
	defer srcFile.Close()

	// Open file2 for writing (without truncating)
	dstFile, err := os.OpenFile("./html/home.jet", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open home.jet: %v", err)
	}
	defer dstFile.Close()

	// Copy content from file1 to file2
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatalf("Failed to copy content: %v", err)
	}

	log.Println("Content copied successfully !!")
}
