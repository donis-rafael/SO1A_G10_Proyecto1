package main

import (
	"context"
	"fmt"
	// Libreria de Google PubSub
	"cloud.google.com/go/pubsub"
)

func publish(msg string) error {
	projectId := "august-edge-306320"
	topicId := "mensajeria"

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		fmt.Println("Error :(")
		return fmt.Errorf("Error al conectarse %v", err)
	}

	t := client.Topic(topicId)

	result := t.Publish(ctx, &pubsub.Message {Data: []byte(msg)})

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return fmt.Errorf("Error: %v",err)
	}

	fmt.Println("Publicando: %v", id)
	return nil
}

func main(){
	fmt.Println("Iniciando envio...")

	publish("Hola mundo desde Go..XD")
}