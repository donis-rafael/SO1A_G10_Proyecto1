package main

import (
	"context"
	"fmt"
	//"io"
	"sync"
	"log"

	"cloud.google.com/go/pubsub"
)

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
			fmt.Println("Got message: %q\n", string(msg.Data))
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