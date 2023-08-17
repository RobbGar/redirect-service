package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/go-redis/redis/v8"
)


const (
	kafkaTopic = "statistics"
)

var (
	serviceURL = os.Getenv("SERVICE_URL")
    kafkaBroker = os.Getenv("KAFKA_BROKER")
	redisURL = os.Getenv("REDIS_URL")
	rdb = redis.NewClient(&redis.Options{
        Addr:     redisURL, 
        Password: "",               
        DB:       0,
    })
    ctx = context.Background()
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println(pong)


    cachedURL, err := rdb.Get(ctx, id).Result()
    
    if err != redis.Nil {  
        if err != nil {
            http.Error(w, "Error reading from cache", http.StatusInternalServerError)
            return
        }
		log.Printf("Read URL from Cache: %s", cachedURL)
        http.Redirect(w, r, cachedURL, http.StatusSeeOther)
        return
    }

	resp, err := http.Get(serviceURL + id)
	if err != nil {
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "URL not found, please make sure that the correct URL has been passed", resp.StatusCode)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	redirectURL := string(bodyBytes)
	err = rdb.Set(ctx, id, redirectURL, 0).Err()
	if err != nil {
	log.Printf("Failed to cache the URL in Redis: %v\n", err)
	} else {
	log.Printf("Put URL: %s in chachen\n", redirectURL)
	}

	publishMessage(id)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func publishMessage(id string) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("id-%s", id)),
		Value: []byte(fmt.Sprintf("Link accessed with id: %s at %s", id, time.Now().String())),
	}

	err := w.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v\n", err)
	} else {
		log.Println("Message successfully sent to Kafka!")
	}

	w.Close()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/r/{id}", redirectHandler)

	http.Handle("/", r)
	port := 8081
	fmt.Printf("Listening on :%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
