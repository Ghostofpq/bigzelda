package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

var redisClient *redis.Client

func main() {
    fmt.Println("hello boyyy")
	InitRedisClient()
	r := mux.NewRouter()
	r.HandleFunc("/gorilla", GorillaHandler)
    r.HandleFunc("/{value}", RedirectionHandler)
    r.HandleFunc("/shortlink/{value}", ShortlinkCreationHandler)
    r.HandleFunc("/admin/{value}", MonitoringHandler)
    http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

func InitRedisClient() {
	redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	pong, err := redisClient.Ping().Result()
	if(pong == "PONG"){
		fmt.Println("Redis started")
	} else{
		fmt.Println(err)		
		os.Exit(1)
	}
}

func GorillaHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gorilla!\n"))
	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
}

func RedirectionHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("RedirectionHandler!\n"))
}

func ShortlinkCreationHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("ShortlinkCreationHandler!\n"))
}

func MonitoringHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("MonitoringHandler!\n"))
}