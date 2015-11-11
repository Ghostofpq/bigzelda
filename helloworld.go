package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

var redisClient *redis.Client

func main() {
	log.Info("Starting up!")
	InitRedisClient()
	r := mux.NewRouter()	
	log.Info("Registering GorillaHandler on /gorilla")
	r.HandleFunc("/gorilla", GorillaHandler)
	log.Info("Registering RedirectionHandler on /{value}")
	r.HandleFunc("/{value}", RedirectionHandler)
	log.Info("Registering ShortlinkCreationHandler on /shortlink/{value}")
	r.HandleFunc("/shortlink/{value}", ShortlinkCreationHandler)
	log.Info("Registering MonitoringHandler on /admin/{value}")
	r.HandleFunc("/admin/{value}", MonitoringHandler)
	http.Handle("/", r)
	log.Info("Starting up on :8000")
	http.ListenAndServe(":8000", r)
}

func InitRedisClient() {
	log.Info("Setting up Redis client")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := redisClient.Ping().Result()
	if pong == "PONG" {
		log.Info("Redis client is up and running")
	} else {
		log.WithFields(log.Fields{"err": err}).Fatal("Redis could not start")
		os.Exit(1)
	}
}

func GorillaHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
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
