package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

var redisClient *redis.Client

func main() {
	log.Info("Starting up!")
	// INIT Redis
	InitRedisClient()

	// INIT Server
	r := mux.NewRouter()

	log.Info("Registering RedirectionHandler on /{value}")
	r.HandleFunc("/{value}", RedirectionHandler).
		Methods("GET")

	log.Info("Registering ShortlinkCreationHandler on /shortlink/{value}")
	r.HandleFunc("/shortlink/{value}", ShortlinkCreationHandler).
		Methods("GET")

	log.Info("Registering MonitoringHandler on /admin/{value}")
	r.HandleFunc("/admin/{value}", MonitoringHandler).
		Methods("GET")

	http.Handle("/", r)
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

func RedirectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortlink := vars["value"]
	log.WithFields(log.Fields{"shortlink": shortlink, "origin": ""}).Debug("redirection")
	http.Redirect(w, r, "http://www.google.com/", http.StatusFound)
}

func ShortlinkCreationHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ShortlinkCreationHandler!\n"))
}

func MonitoringHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MonitoringHandler!\n"))
}
