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

	assignedRoles map[string]*Client  // track role to client mapping, keys: "X", "O"
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast: make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		GameState: &GameState {
			GameStarted: false,
			GameOver: false,
			PlayersConnected: 0,
			Board: [9]string{},
			PlayerTurn: "X",
		},
        assignedRoles: make(map[string]*Client),
	}
}

// Assign a connected client a role baased on available roles
func (h *Hub) assignRole(c *Client) {
    if _, taken := h.assignedRoles["X"]; !taken {
        c.role = "X"
        h.assignedRoles["X"] = c
        return
    }
    if _, taken := h.assignedRoles["O"]; !taken {
        c.role = "O"
        h.assignedRoles["O"] = c
        return
    }
    // otherwise, spectator
    c.role = ""
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
					incPlayerCount := true
					h.GameState.updatePlayerCount(incPlayerCount)
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

				if client.role == "X" || client.role == "O" {
					// free disconnected player's role
					delete(h.assignedRoles, client.role)
					h.GameState.updatePlayerCount(false)
					client.role = ""
				}

				// reset game if player count < 2 and game started
				if h.GameState.PlayersConnected < 2 && h.GameState.GameStarted {
					h.GameState.GameStarted = false
					h.GameState.GameOver = false
					h.GameState.Board = [9]string{}
					h.GameState.PlayerTurn = "X"
				}

				// send game update to clients due to disconnection
				confirm := Message{
                    Type: Game,
					GameState: h.GameState,
                }

				out, _ := json.Marshal(confirm)
				h.broadcast <- out
				
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
