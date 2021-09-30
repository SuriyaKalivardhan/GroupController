package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AllocateRequest struct {
	Id     string `json:"id"`
	Target int    `json:"target"`
}

func main() {
	result := getHello("Starting")
	log.Println(result)

	redisHost := "localhost:6379"
	redisPasswd := ""
	client := initRedis(redisHost, redisPasswd)
	go monitorRedisKeys(client)

	http.HandleFunc("/healthcheck", handleHealthcheck)
	http.HandleFunc("/allocate", handleAllocate)
	http.ListenAndServe(":5001", nil)
}

func handleHealthcheck(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Recived the request: %v", request.Method)
	writer.Write([]byte("Okay"))
}

func handleAllocate(writer http.ResponseWriter, request *http.Request) {
	log.Println("Recieved the request")
	var allocateRequest AllocateRequest
	err := json.NewDecoder(request.Body).Decode(&allocateRequest)
	if err != nil {
		log.Printf("Unexpected response while parsing the request %v", err)
	}
	result := fmt.Sprintf("%s acked with the target of %x", allocateRequest.Id, allocateRequest.Target)
	writer.Write([]byte(result))
}

func getHello(name string) string {
	return "Hello " + name
}

func getSampleRquest() *AllocateRequest {
	return &AllocateRequest{
		Id:     "suriya",
		Target: 34,
	}
}

func copyRequest(request *AllocateRequest) AllocateRequest {
	return AllocateRequest{
		Id:     request.Id,
		Target: request.Target,
	}
}
