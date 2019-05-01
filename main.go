package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adamhe17/board/connect6"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	players    [2]string
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{
				Type:    mtGameState,
				Content: strconv.Itoa(conn.player) + "," + game.GetState(),
			})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, conn *Client) {
	conn.send <- message
}

func (manager *ClientManager) sendAll(message []byte) {
	for conn := range manager.clients {
		conn.send <- message
	}
}

func (manager *ClientManager) sendToPlayer(message []byte, player int) {
	for conn := range manager.clients {
		if conn.id == players[player-1] {
			conn.send <- message
		}
	}
}

func (manager *ClientManager) addPlayer(client *Client) {
	uid := client.id
	if players[0] == "" || players[0] == uid {
		players[0] = uid
		client.player = 1
	} else if players[1] == "" || players[1] == uid {
		players[1] = uid
		client.player = 2
	} else {
		client.player = -1
	}
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
	player int
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}

		moveString := string(message)
		move := connect6.ParseMove(moveString)
		player, gameErr := game.MakeMove(move)
		var messageJSON []byte
		if gameErr != nil {
			messageJSON, _ = json.Marshal(&Message{
				Sender:  c.id,
				Type:    mtError,
				Content: gameErr.Error(),
			})
		} else {
			messageJSON, _ = json.Marshal(&Message{
				Sender:  c.id,
				Type:    mtMove,
				Content: moveString + "," + strconv.Itoa(player),
			})
			manager.broadcast <- messageJSON

			if game.Finished {
				messageJSON, _ = json.Marshal(&Message{
					Sender:  c.id,
					Type:    mtFinished,
					Content: strconv.Itoa(game.CurrentPlayer),
				})
			} else {
				messageJSON, _ = json.Marshal(&Message{
					Sender:  c.id,
					Type:    mtCurrentPlayer,
					Content: strconv.Itoa(game.CurrentPlayer),
				})
			}
		}
		manager.broadcast <- messageJSON
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

type messageType string

type Message struct {
	Sender    string      `json:"sender,omitempty"`
	Recipient string      `json:"recipient,omitempty"`
	Type      messageType `json:"type,omitempty"`
	Content   string      `json:"content,omitempty"`
}

const (
	mtError         messageType = "Error"
	mtGameState     messageType = "GameState"
	mtMove          messageType = "Move"
	mtFinished      messageType = "Finished"
	mtCurrentPlayer messageType = "CurrentPlayer"
)

var (
	manager = ClientManager{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	game    = connect6.Connect6{}
	players [2]string
)

func main() {
	fmt.Println("Starting game")
	game.NewGame()

	fmt.Println("Starting application...")
	go manager.start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	fmt.Println(req.Body)

	uid, _ := uuid.NewRandom()
	client := &Client{
		id:     uid.String(),
		socket: conn,
		send:   make(chan []byte),
	}
	manager.addPlayer(client)

	manager.register <- client

	go client.read()
	go client.write()
}
