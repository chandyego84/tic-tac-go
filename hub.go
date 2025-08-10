// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

    GameState  *GameState
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		GameState: &GameState {
			GameStarted: false,
			GameOver: false,
			PlayersConnected: 0,
			Board: [9]string{},
			PlayerTurn: "",
		},
	}
}

/*
(Goroutine) Waits on multiple communication operations (channels)
*/
func (h *Hub) run() {
	for {
		select {
			case client := <-h.register:
				h.clients[client] = true
				if (client.role != "") {
					addPlayer := true
					client.updatePlayerCount(addPlayer)
				}
				fmt.Printf("New client connected. Total clients: %d\n", len(h.clients))

			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
				
			case message := <-h.broadcast:
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
		}
	}
}
