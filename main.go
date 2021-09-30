package main

import (
	"log"
	"net/http"
)

type AllocateRequest struct {
	id     string
	target int
}

func main() {
	result := getHello("Starting")
	log.Println(result)

	http.HandleFunc("/healthcheck", handleHealthcheck)
	http.ListenAndServe(":5001", nil)
}

func handleHealthcheck(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Recived the request: %v", request.Method)
	writer.Write([]byte("Okay"))
}

func getHello(name string) string {
	return "Hello " + name
}

func getSampleRquest() *AllocateRequest {
	return &AllocateRequest{
		id:     "suriya",
		target: 34,
	}
}

func copyRequest(request *AllocateRequest) AllocateRequest {
	return AllocateRequest{
		id:     request.id,
		target: request.target,
	}
}
