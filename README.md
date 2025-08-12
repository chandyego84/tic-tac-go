# Noob-Tac-Go - Multiplayer Tic Tac Toe on AKS

#### *I did this to learn and practice Go + AKS concepts in ~1 week. The code is not very pretty.*

This is the next-gen Tic Tac Toe game server written in Go using WebSockets. It is containerized with Docker and deployed on Azure Kubernetes Service (AKS).  
The goal is for Noob-Tac-Go/Tic-Tac-Go the #1 browser game. 

## Features

- Real-time gameplay and chat messaaging with WebSocket communication  
- Player role assignment, turn management, and ability to spectate live gameplay  
- Containerized for cloud-native deployment  
- Exposed via a Kubernetes LoadBalancer service with a static public IP and DNS

## Architecture

- **Backend:** Go WebSocket server 
- **Container:** Multi-stage Docker image for lean deployment (actually, it is not lean right now. It's like 1.39 GB)
- **Orchestration:** Kubernetes (AKS) cluster with 2 nodes
- **Networking:** Azure LoadBalancer service with static public IP and DNS

## Role Assignment / Who Gets to Play
There is one server. First-come-first-serve basis. If you're one of the first two to enter, you can play. Otherwise, you are just a spectator, but you can still participate in chat. If a player disconnects, the game is restarted, and the next connected client replaces that player--yes, this is a terrible solution, but I wanted something up and running ASAP.

## Access the App
Open a web browser and navigate to:  
[http://tictacgoapp.westus2.cloudapp.azure.com/](http://tictacgoapp.westus2.cloudapp.azure.com/)
* Get a friend to join you (or open another tab), and pray you are both players because there's just so many people trying to access the server to be the best Noob-Tac-Go player.

## OR Run Locally
### Prerequisites
- Go 1.24 or newer installed  
- Docker installed (optional, for containerized run)  
---
1. Clone the repository and navigate into it.
   ```bash
   git clone https://github.com/chandyego84/tic-tac-go.git
   cd tic-tac-go
    ```
2. Download dependencies and build the binary.
    ```bash
    go mod download
    go build -o tictacgo
    ```
3. Run the app.
    ```bash
    ./tictacgo
    ```
4. The app listens on port 8080 by default.
    ```bash
    http://localhost:8080
    ```

## OR Run With Docker
1. Build the image
    ```bash
    docker build -t tictacgo:latest -f Dockerfile .
    ```
2. Run the container
    ```bash
    docker run -p 8080:808 tictacgo:latest
    ```
3. Open your browser (or multiple to play/chat with yourself) and visit:
    ```bash
    http://localhost:8080
    ```