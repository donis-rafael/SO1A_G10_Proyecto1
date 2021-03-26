package main

import (
	"context"
	"fmt"
	"os"
	"log"
	// Para oir a peticiones GET Y POST
    "net/http"
	// Enviar datos en json
	"encoding/json"

	// Leer variables de entorno
	"github.com/joho/godotenv"
	// Libreria de Google PubSub
	"cloud.google.com/go/pubsub"
)

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")
	
	if err != nil {
	  log.Fatalf("Error cargando las variables de entorno")
	}
	
	return os.Getenv(key)
}

func publish(msg string) error {
	projectID := goDotEnvVariable("PROJECT_ID")
	topicID := goDotEnvVariable("TOPIC_ID")

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println("Error encontrado")
		return fmt.Errorf("Error al conectarse %v", err)
	}

	t := client.Topic(topicID)

	result := t.Publish(ctx, &pubsub.Message {Data: []byte(msg), })

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return fmt.Errorf("Error encontrado: %v",err)
	}

	fmt.Println("Published a message; msg ID: %v\n", id)
	return nil
}

type Message struct {
	//Msg  string
	name string
	location string
	age string
	infectedtype string
	state string
	origen string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func http_server(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    switch r.Method {
		case "GET":     
			http.ServeFile(w, r, "form.html")
			/*w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"message\": \"ok gRPC\"}"))*/
			return

		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			// Obtener el nombre enviado desde la forma
			//name := r.FormValue("name")
			// Obtener el mensaje enviado desde la forma
			msg := r.FormValue("msg")
			fmt.Println(string(msg))

			/*var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			failOnError(err, "Parsing JSON")
			body["origen"] = "PubSub"*/

			
			message, err := json.Marshal(msg)
			//message, err := json.Marshal(Message{name: "asdf" , location:"loc", age: "23", infectedtype: "asd", state: "as", origen: "w4" })
			//message, err := json.Marshal(body)
			// Existio un error generando el objeto JSON
			if err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			fmt.Println(string(message))

			publish(string(message))

			fmt.Fprintf(w, "¡Mensaje Publicado!\n")
			fmt.Fprintf(w, "Message = %s\n", message)
			fmt.Fprintln(w, string(message))
		
		default:
			fmt.Fprintf(w, "Metodo %s no soportado \n", r.Method)
			return
    }
}


func main(){
	fmt.Println("Server Google PubSub iniciado")

	http.HandleFunc("/", http_server)

	http_port := ":" + goDotEnvVariable("PORT")
	
    if err := http.ListenAndServe(http_port, nil); err != nil {
        log.Fatal(err)
    }
}