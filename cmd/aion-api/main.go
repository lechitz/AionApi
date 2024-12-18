package main

import (
	"fmt"
	"github.com/lechitz/AionApi/src/config"
	"github.com/lechitz/AionApi/src/router"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()
	r := router.GenareteRouter()

	fmt.Println("Starting Aion API...")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
