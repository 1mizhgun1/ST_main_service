package main

import (
	api "github.com/1mizhgun1/ST_main_service/api"
)

const host = "127.0.0.1:8080"

func main() {
	api.StartServer(host)
}
