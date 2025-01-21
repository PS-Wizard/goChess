package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/notnil/chess"
)

type Server struct {
	Port         string
	Rooms        map[string]*Room 
	Mutex        sync.Mutex       
}

type Room struct {
	Game        *chess.Game       
	Connections []*websocket.Conn 
	CurrentTurn chess.Color       
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func NewServer(port string) *Server {
	return &Server{
		Port:  port,
		Rooms: make(map[string]*Room),
	}
}

func (s *Server) StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.handleWebSocket(w, r)
	})

	fmt.Printf("Starting server at %s\n", s.Port)
	err := http.ListenAndServe(s.Port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}


func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "default"
	}

	s.Mutex.Lock()
	room, exists := s.Rooms[roomName]
	if !exists {
		room = &Room{
			Game:        chess.NewGame(), 
			Connections: []*websocket.Conn{},
		}
		s.Rooms[roomName] = room
		room.CurrentTurn = chess.White 
	}
	room.Connections = append(room.Connections, conn)
	s.Mutex.Unlock()

	fmt.Printf("User joined room: %s\n", roomName)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		s.handleMove(roomName, msg, conn)
	}
}



func (s *Server) handleMove(roomName string, msg []byte, conn *websocket.Conn) {
	var moveData struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}

	err := json.Unmarshal(msg, &moveData)
	if err != nil {
		fmt.Println("Invalid message format:", err)
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	room, exists := s.Rooms[roomName]
	if !exists {
		fmt.Println("Room does not exist:", roomName)
		return
	}

	if room.Game.Position().Turn() != room.CurrentTurn {
		response := map[string]string{
			"type":    "error",
			"message": "It's not your turn!",
		}
		err := conn.WriteJSON(response)
		if err != nil {
			fmt.Println("Error sending turn error:", err)
		}
		return
	}

	var selectedMove *chess.Move
	for _, move := range room.Game.ValidMoves() {
		if move.S1().String() == moveData.Source && move.S2().String() == moveData.Target {
			selectedMove = move
			break
		}
	}

	if selectedMove == nil {
		response := map[string]string{
			"type":    "error",
			"message": "Invalid move! Please try again.",
		}
		err := conn.WriteJSON(response)
		if err != nil {
			fmt.Println("Error sending error message:", err)
		}
		return
	}

	room.Game.Move(selectedMove)

	outcome := room.Game.Outcome()

	var moveResponse map[string]string
	var gameOverResponse map[string]string

	if outcome == chess.WhiteWon || outcome == chess.BlackWon {
		moveResponse = map[string]string{
			"type":    "move",
			"newPos":  room.Game.Position().String(),
			"message": fmt.Sprintf("Move: %s -> %s (%s won!)", moveData.Source, moveData.Target, outcome),
		}
		gameOverResponse = map[string]string{
			"type":    "gameOver",
			"message": fmt.Sprintf("Game over! %s won!", outcome),
		}
		s.broadcastMessage(room, moveResponse)
		s.broadcastMessage(room, gameOverResponse)
	} else if outcome == chess.Draw {
		moveResponse = map[string]string{
			"type":    "move",
			"newPos":  room.Game.Position().String(),
			"message": fmt.Sprintf("Move: %s -> %s (Draw!)", moveData.Source, moveData.Target),
		}
		gameOverResponse = map[string]string{
			"type":    "gameOver",
			"message": "Game over! It's a draw!",
		}
		s.broadcastMessage(room, moveResponse)
	} else {
		fen := room.Game.Position().String()
		moveResponse = map[string]string{
			"type":    "move",
			"newPos":  fen,
			"message": fmt.Sprintf("Move: %s -> %s", moveData.Source, moveData.Target),
		}
		s.broadcastMessage(room, moveResponse)
	}

	room.CurrentTurn = room.Game.Position().Turn()
}

func (s *Server) broadcastMessage(room *Room, data interface{}) {
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	for _, conn := range room.Connections {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error broadcasting message:", err)
		}
	}
}
