package main

import (
	"github.com/quangtran666/godot-golangudp-ping-pong/internal/network"
	"log"
)

func main() {
	server := network.NewServer()
	if err := server.Start("127.0.0.1:8080"); err != nil {
		log.Printf("Error starting server: %v", err)
	}
}
