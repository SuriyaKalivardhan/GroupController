package main

import "log"

type AllocateRequest struct {
	id     string
	target int
}

func main() {
	result := getHello("Agazhi")
	log.Println(result)
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
