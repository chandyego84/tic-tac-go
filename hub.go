// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
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
			PlayerTurn: "X",
		},
	}
}

// Assign a connected client a role based on number of players connected
// First two clients to connect are players
func (h *Hub) assignRole(c *Client) {
	switch h.GameState.PlayersConnected {
		case 0: {
			c.role = "X"
		}
		case 1: {
			c.role = "O"
		} 
		default: c.role = ""
	}
}

func (c *Client) updateClientRegisterStatus(isRegistered bool) {
	c.registered = isRegistered
}

/*
(Goroutine) Waits on multiple communication operations (channels)
*/
func (h *Hub) run() {
	for {
		select {
			case client := <-h.register:
				h.clients[client] = true
				client.updateClientRegisterStatus(true)
				
				h.assignRole(client)
				if client.role != "" {
					addToPlayerCount := true
					h.GameState.updatePlayerCount(addToPlayerCount)
				}

				// send role confirmation to client along with game state
				confirm := Message{
                    Type: Connection,
					Username: client.username,
					Role: client.role,
					GameState: h.GameState,
                }

                out, _ := json.Marshal(confirm)
				client.send <- out

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
