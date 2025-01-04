package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port string
}

func NewServer(port string) *Server {
	return &Server{
		Port: port,
	}
}

func (s *Server) StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
	})
	fmt.Printf("Starting server at %s\n", s.Port)
	err := http.ListenAndServe(s.Port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

