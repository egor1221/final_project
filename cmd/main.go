package main

import (
	"log"
	"todo/server"
)

func main() {

	err := server.StartServer()

	if err != nil {
		log.Fatalf(err.Error())
	}
}
