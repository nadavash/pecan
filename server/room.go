package main

import (
	"log"
)

// Room maintains a set of active controller clients and one game client to receive controller
// inputs.
type Room struct {

	// Client that receives controller inputs.
	gameClient *Client

	// Clients that send controller input to the gameClient.
	controllerClients map[*Client]bool

	controllerInput chan []byte

	unregister chan *Client
}

// NewRoom creates a new room for the given game Client.
func NewRoom(gameClient *Client) *Room {
	r := &Room{
		gameClient:        gameClient,
		controllerClients: make(map[*Client]bool),
		controllerInput:   make(chan []byte),
		unregister:        make(chan *Client),
	}
	go gameClient.readPump(r.controllerInput, r.unregister)
	go gameClient.writePump()
	return r
}

// RegisterControllerClient registers a controller client for this room.
func (r *Room) RegisterControllerClient(client *Client) bool {
	r.controllerClients[client] = true
	go client.readPump(r.controllerInput, r.unregister)
	go client.writePump()
	return true
}

// Run runs the main logic for the Room.
func (r *Room) Run() {
	defer func() {
		log.Println("Shutting down room.")
		r.shutdown()
	}()

	for {
		select {
		case client := <-r.unregister:
			if r.gameClient == client {
				log.Println("Game client disconnected.")
				return
			} else if _, ok := r.controllerClients[client]; ok {
				log.Println("Controller client disconnected.")
				delete(r.controllerClients, client)
				client.Close()
			}
		case message := <-r.controllerInput:
			if err := r.gameClient.sendMessage(message); err != nil {
				log.Printf("error: %v\n", err)
				return
			}
		}
	}
}

// Shutdown shuts down the room, disconnecting all of the clients.
func (r *Room) shutdown() {
	r.gameClient.Close()
	for client := range r.controllerClients {
		client.Close()
	}
}
