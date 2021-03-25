package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

func elemento(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	failOnError(err, "Parsing JSON")
	body["origen"] = "Nats"
	data, err := json.Marshal(body)

	nc, err := nats.Connect("nats://servern:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// type Person struct {
	// 	Name         string `json:"name"`
	// 	Location     string `json:"location"`
	// 	Age          int    `json:"age"`
	// 	Infectedtype string `json:"infectedtype"`
	// 	State        string `json:"state"`
	// }

	// body, errs := ioutil.ReadAll(r.Body)
	// if errs != nil {
	// 	log.Printf("Error reading body: %v", errs)
	// 	http.Error(w, "can't read body", http.StatusBadRequest)
	// 	return
	// }

	// myBody := ioutil.NopCloser(bytes.NewBuffer(body))

	// var p3 Person
	// dec := json.NewDecoder(myBody)
	// dec.DisallowUnknownFields()
	// erre := dec.Decode(&p3)

	newData := string(data)
	nc.Publish("colaid", []byte(newData))
	log.Println("pub finish")

	nc.Flush()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(newData))
}

func handleRequests() {
	http.HandleFunc("/", elemento)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	handleRequests()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
