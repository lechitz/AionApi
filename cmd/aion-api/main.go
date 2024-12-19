package main

import (
	"fmt"
	"github.com/lechitz/AionApi/src/config"
	"github.com/lechitz/AionApi/src/router"
	"log"
	"net/http"
)

//func init() {
//	key := make([]byte, 64)
//
//	if _, err := rand.Read(key); err != nil {
//		log.Fatal(err)
//	}
//
//	stringBase64 := base64.StdEncoding.EncodeToString(key)
//
//	fmt.Println(stringBase64)
//}

func main() {
	config.LoadEnv()
	r := router.GenareteRouter()

	fmt.Println("Starting Aion API...")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
