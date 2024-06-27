package main

import (
	"desafio-client-server/client"
	"desafio-client-server/server"
	"log"
	"net/http"
	"time"
)

var ServerURL = "http://localhost:8080/cotacao"

func main() {
	go func() {
		log.Println("Starting server...")
		server.StartServer()
	}()

	waitForServer(ServerURL)

	client.StartClient()
}

func waitForServer(url string) {
	for {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Println("Server is up and running.")
			return
		}
		log.Println("Waiting for the server to be ready...")
		time.Sleep(2 * time.Second)
	}
}
