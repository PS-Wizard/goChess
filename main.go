package main

import "goChess/server"

func main() {
    server := server.NewServer(":8080")
    server.StartServer();
}

