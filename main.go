package main

import (
	"fmt"
	"log"
	"net/http"
	"redirect-service/api"
	"redirect-service/redis"
)

func main() {
	redis.InitializeRedis()
	api.SetupRoutes()

	port := 8081
	fmt.Printf("Listening on :%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
