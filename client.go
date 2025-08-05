// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// name of peer
	clientName string

	// Buffered channel of outbound messages to clients.
	send chan []byte
}

// Message wraps the user and the message data.
type Message struct {
	user string

	message string
}

// readPump reads messages from the client and sends them to the hub
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := newMessage(c.clientName, string(message))
		outMessage := bytes.TrimSpace([]byte(msg.user + ": " + msg.message))
		c.hub.broadcast <- outMessage
	}
}

// writePump writes messages from the hub to the client (sends outbound messages)
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod) // sends periodic ping messages to keep the connection alive

	defer func() {
		// Stop the ticker and close the WebSocket connection when this function exits
		ticker.Stop()
		c.conn.Close()
	}()

	// Loop forever, waiting to either:
	// - receive a message to send from c.send
	// - or send a periodic ping to keep the connection alive
	for {
		select {

		// Case 1: message ready to be sent to client
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// The `c.send` channel was closed (e.g., hub unregistered the client)
				// Send a close message to the WebSocket client
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Start writing a new WebSocket message of type "TextMessage"
			// This returns a writer that allows streaming large messages
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				// If we can't get a writer (maybe connection is closed), stop
				return
			}

			// Write the current message into the writer
			w.Write(message)

			// Also write any **queued messages** that came in while we were writing
			// This helps batch multiple messages into a single WebSocket frame
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			// Finish writing and flush everything to the network
			if err := w.Close(); err != nil {
				return
			}

		// Case 2: time to send a ping message to keep the connection alive
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			// Send a Ping control message (used to detect dropped connections)
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Creates a new message with user and message content.
// This is a utility function to encapsulate message creation.
func newMessage(user string, message string) *Message {
	return &Message{
		user:    user,
		message: message,
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client_num := strconv.Itoa(rand.Int()) // generate random user number for client name
	client := &Client{hub: hub, conn: conn, clientName: "user" + client_num, send: make(chan []byte, 256)}
	client.hub.register <- client // send client to hub's register channel

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
