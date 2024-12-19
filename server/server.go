package server

import (
	"fmt"
	"goChess/board"
	"net/http"
)

type Server struct {
	Port   string
	router *http.ServeMux
	Board  *board.Board
}

func NewServer(port string) *Server {
	return &Server{
		Port:   port,
		router: http.NewServeMux(),
	}
}
func (s *Server) StartServer() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Board = board.NewBoard()
		http.ServeFile(w, r, "./ui/index.html")
	})
	s.router.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		peice := r.URL.Query().Get("piece")
        from := r.URL.Query().Get("from")
        to := r.URL.Query().Get("to")
        old := r.URL.Query().Get("oldState")
		new := r.URL.Query().Get("newState")
		fmt.Println("Got Request", old, new)
        fmt.Fprintf(w, "got request, To Move: %s , from: %s , to: %s", peice, from, to)
        // TODO:
        //     - Check If Move is valid,
        //     - If It Is Its Bueno, Update Internal state
        //     - If Not Return Html Snippet for chess.js with position of old
	})
	fmt.Printf("Starting a server at %s", s.Port)
	http.ListenAndServe(s.Port, s.router)
}
