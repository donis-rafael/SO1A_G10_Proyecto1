package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"log"
	"bytes"
	"net/http"
	"cloud.google.com/go/pubsub"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMongo(d string){
	postBody := []byte(string(d))
	req, err := http.Post("http://34.121.234.71:3000/subscribers", "application/json", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	failOnError(err, "POST new document")
	defer req.Body.Close()

	//Read the response body
	newBody, err := ioutil.ReadAll(req.Body)
	failOnError(err, "Reading response from HTTP POST")
	sb := string(newBody)
	log.Printf(sb)
}

func main(){
	fmt.Println("A la espera de mensajes...")

	projectID := "august-edge-306320"
	subID := "mensaje"

	ctx := context.Background()
	
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
			//return fmt.Errorf("pubsub.NewClient: %v", err)
			log.Fatal(err)
	}

	// Consume 10 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
			mu.Lock()
			defer mu.Unlock()
			fmt.Println("Mensaje en suscriptor recevido: %q\n", string(msg.Data))
			
			sendMongo(string(msg.Data))

			msg.Ack()
			received++
			if received == 100 {
					cancel()
			}
	})
	if err != nil {
			//return fmt.Errorf("Receive: %v", err)
			log.Fatal(err)
	}

}
