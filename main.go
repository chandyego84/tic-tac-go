package main
import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

/*
* Upgrades an HTTP connection to a websocket
*/
var upgrader = websocket.Upgrader {
	// allow all connections (in prod, should validate origin to avoid cross-site websocket hijacking)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// upgrade GET request to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if (err != nil) {
		fmt.Println(err.Error())
		return
	}
	defer ws.Close() // stupid language

	for { // read/write loop -- read messages from client and echo them back
		// read msgs from browser
		_, msg, err := ws.ReadMessage()
		if (err != nil) {
			fmt.Println("read error: ", err.Error())
			break
		}
		fmt.Printf("Received: %s\n", msg)

		// write msg back to browser
        if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil { 
            fmt.Println("write error:", err)
            break
        }
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Websocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if (err != nil) {
		fmt.Println("ListenAndServe: ", err.Error())
	}
}