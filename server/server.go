package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var flagAddr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var gameRoom *Room = nil

func serveWebsocketGame(w http.ResponseWriter, r *http.Request) {
	log.Println("New game connected")
	if gameRoom != nil {
		http.Error(w, "Game client already exists.", http.StatusConflict)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading WS connection:", err)
		return
	}
	gameClient := NewClient(conn)
	gameRoom = NewRoom(gameClient)
	gameRoom.Run()
	gameRoom = nil
}

func serveWebsocketController(w http.ResponseWriter, r *http.Request) {
	log.Println("New controller connected")
	if gameRoom == nil {
		http.Error(w, "Game client does not exist.", http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading WS connection:", err)
		return
	}
	controllerClient := NewClient(conn)
	gameRoom.RegisterControllerClient(controllerClient)
}

func main() {
	flag.Parse()

	http.HandleFunc("/ws/game", serveWebsocketGame)
	http.HandleFunc("/ws/controller", serveWebsocketController)

	err := http.ListenAndServe(*flagAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
