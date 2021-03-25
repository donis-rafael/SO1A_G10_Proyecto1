package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var received int64

func main() {
	received = 0

	go subscriptor()
	go check()

	for {

	}
}

func subscriptor() {

	nc, err := nats.Connect("nats://servern:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	//Nos suscribimos para escuchar mensajes
	nc.Subscribe("colaid", func(msg *nats.Msg) {
		log.Printf("%s: %s", msg.Subject, msg.Data)
		received++

		log.Printf("Mensaje Recibido: %s", msg.Data)

		postBody := []byte(string(msg.Data))

		req, err := http.Post("http://34.121.234.71:3000/subscribers", "application/json", bytes.NewBuffer(postBody))
		req.Header.Set("Content-Type", "application/json")
		failOnError(err, "Recibido")
		defer req.Body.Close()

		//Read the response body
		newBody, err := ioutil.ReadAll(req.Body)
		failOnError(err, "HTTP POST")
		sb := string(newBody)
		log.Printf(sb)
		log.Printf("Cancel CTRL+C")

	})
	nc.Flush()

	for {

	}
}

func check() {
	for {
		fmt.Println("-----------------------")
		fmt.Println("still running")
		fmt.Println("received", received)
		fmt.Println("-----------------------")
		time.Sleep(time.Second * 2)
	}
}
