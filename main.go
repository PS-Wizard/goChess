package main

import "goChess/server"

func main() {
    server := server.NewServer("0.0.0.0:8080")
    server.StartServer();
}
