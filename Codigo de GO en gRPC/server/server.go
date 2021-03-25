// Paquete principal, acá iniciará la ejecución
package main

// Importar dependencias, notar que estamos en un módulo llamado grpctuiter
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"net/http"

	"os"

	"log"

	"servergrpc/greet.pb/greetpb"

	"google.golang.org/grpc"
)

// Iniciar una estructura que posteriormente gRPC utilizará para realizar un server
type server struct{}

// Función que será llamada desde el cliente
// Debemos pasarle un contexto donde se ejecutara la funcion
// Y utilizar las clases que fueron generadas por nuestro proto file
// Retornara una respuesta como la definimos en nuestro protofile o un error
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf(">> SERVER: Función Greet llamada con éxito. Datos: %v\n", req)

	// Todos los datos podemos obtenerlos desde req
	// Tendra la misma estructura que definimos en el protofile
	// Para ello utilizamos en este caso el GetGreeting
	Name := req.GetGreeting().GetName()
	Location := req.GetGreeting().GetLocation()
	Age := req.GetGreeting().GetAge()
	Infectedtype := req.GetGreeting().GetInfectedtype()
	State := req.GetGreeting().GetState()

	result := Name + " - " + Location + " - " + Age + " - " + Infectedtype + " - " + State

	fmt.Printf(">> SERVER: %s\n", result)
	// Creamos un nuevo objeto GreetResponse definido en el protofile

	jsonData := map[string]string{"name": Name, "location": Location, "age": Age, "infectedtype": Infectedtype, "state": State}
	jsonValue, _ := json.Marshal(jsonData)
	//client := &http.Client{}
	request, err := http.Post("http://node:3000/subscribers", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Print(err.Error())
	}
	//result.Header.Add("Accept", "application/json")
	//result.Header.Add("Content-Type", "application/json")

	//resp, err := client.Do(request)
	//request, err = http.Post("https://localhost/subscribers", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(request.Body)
		fmt.Println(string(data))
	}
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

// Funcion principal
func main() {

	// Leer el host de las variables del ambiente
	host := os.Getenv("HOST")
	fmt.Println(">> SERVER: Iniciando en ", host)

	// Primero abrir un puerto para poder escuchar
	// Lo abrimos en este puerto arbitrario
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf(">> SERVER: Error inicializando el servidor: %v", err)
	}

	fmt.Println(">> SERVER: Empezando server gRPC")

	// Ahora si podemos iniciar un server de gRPC
	s := grpc.NewServer()

	// Registrar el servicio utilizando el codigo que nos genero el protofile
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Println(">> SERVER: Escuchando servicio...")
	// Iniciar a servir el servidor, si hay un error salirse
	if err := s.Serve(lis); err != nil {
		log.Fatalf(">> SERVER: Error inicializando el listener: %v", err)
	}
}
