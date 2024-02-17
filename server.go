package main

import (
	"log"

	api "github.com/1mizhgun1/ST_main_service/api"
)

func main() {
	log.Println("App start")
	api.StartServer()
	log.Println("App stop")
}
