package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var client = http.DefaultClient
var url = "http://localhost:8080/cotacao"

type BidResponse struct {
	Bid string `json:"bid"`
}

func StartClient() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var bidResponse BidResponse
	err = json.Unmarshal(body, &bidResponse)
	if err != nil {
		log.Println(err)
		return
	}

	err = os.WriteFile("cotacao.txt", []byte("Dólar: "+bidResponse.Bid), 0644)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Bid value: %s\n", bidResponse)
}
