// Paquete principal, acá iniciará la ejecución
package main

// Importar dependencias, notar que estamos en un módulo llamado tuiterclient
import (
	"context"

	"fmt"
	"io/ioutil"
	"log"
	"os"

	"net/http"

	"clientgrpc/greet.pb/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

type Person struct {
	name         string
	location     string
	age          int
	infectedtype string
	state        string
}

// Funcion que realiza una llamada unaria
func sendMessage(name string, location string, age string, infectedtype string, state string) {
	server_host := os.Getenv("SERVER_HOST")

	fmt.Println(">> CLIENT: Iniciando cliente")
	fmt.Println(">> CLIENT: Iniciando conexion con el servidor gRPC ", server_host)

	// Crear una conexion con el servidor (que esta corriendo en el puerto 50051)
	// grpc.WithInsecure nos permite realizar una conexion sin tener que utilizar SSL
	cc, err := grpc.Dial(server_host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(">> CLIENT: Error inicializando la conexion con el server %v", err)
	}

	// Defer realiza una axion al final de la ejecucion (en este caso, desconectar la conexion)
	defer cc.Close()

	// Iniciar un servicio NewGreetServiceClient obtenido del codigo que genero el protofile
	// Esto crea un cliente con el cual podemos escuchar
	// Le enviamos como parametro el Dial de gRPC
	c := greetpb.NewGreetServiceClient(cc)

	fmt.Println(">> CLIENT: Iniciando llamada a Unary RPC")

	// Crear una llamada de GreetRequest
	// Este codigo lo obtenemos desde el archivo que generamos con protofile
	req := &greetpb.GreetRequest{
		// Enviar un Greeting
		// Esta estructura la obtenemos desde el archivo que generamos con protofile
		Greeting: &greetpb.Greeting{
			Name:         name,
			Location:     location,
			Age:          age,
			Infectedtype: infectedtype,
			State:        state,
		},
	}

	fmt.Println(">> CLIENT: Enviando datos al server")
	// Iniciar un greet, en background con la peticion que estamos realizando
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf(">> CLIENT: Error realizando la peticion %v", err)
	}

	fmt.Println(">> CLIENT: El servidor nos respondio con el siguiente mensaje: ", res.Result)
}

// Creamos un server sencillo que unicamente acepte peticiones GET y POST a '/'
func http_server(w http.ResponseWriter, r *http.Request) {
	instance_name := os.Getenv("NAME")
	fmt.Println(">> CLIENT: Manejando peticion HTTP CLIENTE: ", instance_name)
	// Comprobamos que el path sea exactamente '/' sin parámetros

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// // Comprobamos el tipo de peticion HTTP
	// // Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// // response body. A request body larger than that will now result in
	// // Decode() returning a "http: request body too large" error.

	// if r.Header.Get("Content-Type") != "" {
	// 	value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
	// 	if value != "application/json" {
	// 		msg := "Content-Type header is not application/json"
	// 		http.Error(w, msg, http.StatusUnsupportedMediaType)
	// 		fmt.Println("Content-Type header is not application/json")
	// 		return
	// 	}
	// }

	// fmt.Println(">> CLIENT: Recibiendo body: ", r)
	// r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// // Setup the decoder and call the DisallowUnknownFields() method on it.
	// // This will cause Decode() to return a "json: unknown field ..." error
	// // if it encounters any extra unexpected fields in the JSON. Strictly
	// // speaking, it returns an error for "keys which do not match any
	// // non-ignored, exported fields in the destination".
	// dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()

	// var p Person
	// err := dec.Decode(&p)
	// if err != nil {
	// 	var syntaxError *json.SyntaxError
	// 	var unmarshalTypeError *json.UnmarshalTypeError

	// 	switch {
	// 	// Catch any syntax errors in the JSON and send an error message
	// 	// which interpolates the location of the problem to make it
	// 	// easier for the client to fix.
	// 	case errors.As(err, &syntaxError):
	// 		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
	// 		http.Error(w, msg, http.StatusBadRequest)
	// 		fmt.Println(">> CLIENT: Error 10 Recibiendo: ", err.Error())

	// 	// In some circumstances Decode() may also return an
	// 	// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
	// 	// is an open issue regarding this at
	// 	// https://github.com/golang/go/issues/25956.
	// 	case errors.Is(err, io.ErrUnexpectedEOF):
	// 		msg := fmt.Sprintf("Request body contains badly-formed JSON")
	// 		http.Error(w, msg, http.StatusBadRequest)
	// 		fmt.Println(">> CLIENT: Error 11 Recibiendo: ", err.Error())

	// 	// Catch any type errors, like trying to assign a string in the
	// 	// JSON request body to a int field in our Person struct. We can
	// 	// interpolate the relevant field name and position into the error
	// 	// message to make it easier for the client to fix.
	// 	case errors.As(err, &unmarshalTypeError):
	// 		msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	// 		http.Error(w, msg, http.StatusBadRequest)
	// 		fmt.Println(">> CLIENT: Error 12 Recibiendo: ", err.Error())

	// 	// Catch the error caused by extra unexpected fields in the request
	// 	// body. We extract the field name from the error message and
	// 	// interpolate it in our custom error message. There is an open
	// 	// issue at https://github.com/golang/go/issues/29035 regarding
	// 	// turning this into a sentinel error.
	// 	case strings.HasPrefix(err.Error(), "json: unknown field "):
	// 		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
	// 		msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
	// 		http.Error(w, msg, http.StatusBadRequest)
	// 		fmt.Println(">> CLIENT: Error 13 Recibiendo: ", p)

	// 	// An io.EOF error is returned by Decode() if the request body is
	// 	// empty.
	// 	case errors.Is(err, io.EOF):
	// 		msg := "Request body must not be empty"
	// 		http.Error(w, msg, http.StatusBadRequest)
	// 		fmt.Println(">> CLIENT: Error 14 Recibiendo: ", err.Error())

	// 	// Catch the error caused by the request body being too large. Again
	// 	// there is an open issue regarding turning this into a sentinel
	// 	// error at https://github.com/golang/go/issues/30715.
	// 	case err.Error() == "http: request body too large":
	// 		msg := "Request body must not be larger than 1MB"
	// 		http.Error(w, msg, http.StatusRequestEntityTooLarge)
	// 		fmt.Println(">> CLIENT: Error 15 Recibiendo: ", err.Error())
	// 	// Otherwise default to logging the error and sending a 500 Internal
	// 	// Server Error response.
	// 	default:
	// 		log.Println(err.Error())
	// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 		fmt.Println(">> CLIENT: Error 16 Recibiendo: ", err.Error())
	// 	}
	// 	fmt.Println(">> CLIENT: Error Recibiendo: ", err.Error())
	// 	return
	// }

	// // Call decode again, using a pointer to an empty anonymous struct as
	// // the destination. If the request body only contained a single JSON
	// // object this will return an io.EOF error. So if we get anything else,
	// // we know that there is additional data in the request body.
	// err = dec.Decode(&struct{}{})
	// if err != io.EOF {
	// 	msg := "Request body must only contain a single JSON object"
	// 	http.Error(w, msg, http.StatusBadRequest)
	// 	fmt.Println(">> CLIENT: Error 2 Recibiendo: ", err.Error())

	// 	return
	// }
	// fmt.Println(">> CLIENT: Recibiendo: ", p)

	switch r.Method {

	// Devolver una página sencilla con una forma html para enviar un mensaje
	case "GET":
		fmt.Println(">> CLIENT: Devolviendo form.html")
		// Leer y devolver el archivo form.html contenido en la carpeta del proyecto
		http.ServeFile(w, r, "form.html")

	// Publicar un mensaje a Google PubSub
	case "POST":
		fmt.Println(">> CLIENT: Iniciando envio de mensajes")
		//Si existe un error con la forma enviada entonces no seguir
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		fmt.Println(">> BODY: Iniciando  ", body)

		// Obtener el nombre enviado desde la forma
		name := r.FormValue("name")
		// Obtener el mensaje enviado desde la forma
		location := r.FormValue("location")
		age := r.FormValue("age") //strconv.Itoa(p.age)
		infectedtype := r.FormValue("infectedtype")
		state := r.FormValue("state")

		// Publicar el mensaje, convertimos el objeto JSON a String
		sendMessage(name, location, age, infectedtype, state)

		// Enviamos informacion de vuelta, indicando que fue generada la peticion
		fmt.Fprintf(w, "¡Mensaje Publicado!\n")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Location = %s\n", location)
		fmt.Fprintf(w, "Age = %s\n", age)
		fmt.Fprintf(w, "Type = %s\n", infectedtype)
		fmt.Fprintf(w, "State = %s\n", state)

	// Cualquier otro metodo no sera soportado
	default:
		fmt.Fprintf(w, "Metodo %s no soportado \n", r.Method)
		return
	}
}

// Funcion principal
func main() {
	instance_name := os.Getenv("NAME")
	client_host := os.Getenv("CLIENT_HOST")

	fmt.Println(">> -------- CLIENTE ", instance_name, " --------")

	fmt.Println(">> CLIENT: Iniciando servidor http en ", client_host)

	// Asignar la funcion que controlara las llamadas http
	http.HandleFunc("/", http_server)

	// Levantar el server, si existe un error levantandolo hay que apagarlo
	if err := http.ListenAndServe(client_host, nil); err != nil {
		log.Fatal(err)
	}
}
