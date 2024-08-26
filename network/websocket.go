package network

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024, //buffer sizes not max message size limit
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //allow all origins
	},
}

func RunWebsocketExample() {
	go startServer() //separate goroutine ListenAndServe is blocking
	startClient()
}

// starts async server that both listens for incoming connections and sends messages to clients
// in different goroutines
func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/websocket", wsHandler)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

// this takes http connection and upgrades it to websocket connection
func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ws connection from %v", r.RemoteAddr)
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}
	defer ws.Close()
	closed := make(chan any) //channel to unlock handler and call 'defer ws.Close()'

	//async, reads message from client
	go read("SERVER", ws, closed)

	//async, periodically sends message to client
	go write(ws, closed)

	//lock handler until close signal is received either from read or write
	select {
	case <-closed:
		log.Printf("closing connection")
		break
	}
}

// starts client that both listens for incoming messages and sends messages to server
// pretty much the same as server
func startClient() {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/websocket", nil)
	if err != nil {
		log.Fatalf("error dialing server: %v", err)
	}
	defer ws.Close()
	closed := make(chan any) //channel to unlock handler and call 'defer ws.Close()'

	//async, reads message from server
	go read("CLIENT", ws, closed)

	//async, periodically sends message to server
	go write(ws, closed)

	//lock handler until close signal is received either from read or write
	select {
	case <-closed:
		log.Printf("closing connection")
		break
	case <-time.After(15 * time.Second): //exit program after 15 seconds
		break
	}
}

func read(prefix string, ws *websocket.Conn, closed chan any) {
	for {
		//listening for close signal omitted for simplicity
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("connection closed by client")
				closed <- true
				return
			}
			log.Printf("error reading message: %v", err)
			closed <- true
			return
		}
		log.Printf("[%s] received message: %v", prefix, string(msg))
	}
}

func write(ws *websocket.Conn, closed chan any) {
	for {
		//listening for close signal omitted for simplicity
		err := ws.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(rand.Intn(100))))
		time.Sleep(time.Duration(3000+rand.Intn(2000)) * time.Millisecond)
		if err != nil {
			log.Printf("error writing message: %v", err)
			closed <- true
			return
		}
	}
}
