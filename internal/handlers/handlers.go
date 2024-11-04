package handlers

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

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
type WsJsonResponse struct {
	Action       string `json:"action"`
	Message      string `json:"message"`
	Message_Type string `json:"message_type"`
}

func WSEndpoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WSEndpoint 1: ", err)
	}
	log.Println("Client is connected to endpoint !")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to the server</small></em>`

	log.Println("Dispatching ws response...")
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("WSEndpoint 2: ", err)
	}
	log.Println("Dispatching Done !")

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
