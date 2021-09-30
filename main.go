package main

import "log"

func main() {
	result := getHello("Agazhi")
	log.Println(result)
}

func getHello(name string) string {
	return "Hello " + name
}
