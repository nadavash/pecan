package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
)

var (
	flagAddr    = flag.String("addr", ":8080", "http service address")
	flagWebRoot = flag.String("web_root", "web/", "the root directory for web files")

	lanIP    = getLocalIP()
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	gameRoom *Room = nil
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(*flagWebRoot, "index.html.template")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Error parse template for serveIndex: %v", err)
		return
	}

	data := struct{ IPAddress string }{lanIP}
	tmpl.Execute(w, data)
}

func serveWebsocketGame(w http.ResponseWriter, r *http.Request) {
	if gameRoom != nil {
		log.Println("Game client tried connecting, but game already exists!")
		http.Error(w, "Game client already exists.", http.StatusConflict)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading WS connection:", err)
		return
	}
	log.Println("New game connected.")
	gameClient := NewClient(conn)
	gameRoom = NewRoom(gameClient)
	gameRoom.Run()
	gameRoom = nil
}

func serveWebsocketController(w http.ResponseWriter, r *http.Request) {
	if gameRoom == nil {
		log.Println("Controller tried connecting, but game is not available!")
		http.Error(w, "Game client does not exist.", http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading WS connection:", err)
		return
	}
	log.Println("New controller connected.")
	controllerClient := NewClient(conn)
	gameRoom.RegisterControllerClient(controllerClient)
}

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir(*flagWebRoot))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws/game", serveWebsocketGame)
	http.HandleFunc("/ws/controller", serveWebsocketController)

	err := http.ListenAndServe(*flagAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
