package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Ghostofpq/bigzelda/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/Ghostofpq/bigzelda/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/Ghostofpq/bigzelda/Godeps/_workspace/src/gopkg.in/redis.v3"
)

// Redis client
var redisClient *redis.Client

// Shortlink structure
type Shortlink struct {
	Id, Token, Origin string
	CreationTs        int64
	Count             uint8
}

//Shortlink request structure
type AdvancedShortlinkRequest struct {
	Origin, Token string
}

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

	log.Info("Registering AdvancedShortlinkCreationHandler on /shortlink")
	r.HandleFunc("/shortlink", AdvancedShortlinkCreationHandler).
		Methods("POST")

	log.Info("Registering MonitoringHandler on /admin/{value}")
	r.HandleFunc("/admin/{value}", MonitoringHandler).
		Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

// REDIS INIT

func InitRedisClient() {
	log.Info("Setting up Redis client")
	redisClient = redis.NewClient(&redis.Options{
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

// HANDLERS

func RedirectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("RedirectionHandler")

	// Extract request parameters
	vars := mux.Vars(r)
	token := vars["value"]

	// Fetch Shortlink in Redis
	shortlink, err := ReadFromRedis(token)
	if err != nil {
		http.Error(w, "Token not found", 404)
		return
	}

	// Increment
	shortlink.Count++
	UpdateShortlink(shortlink)

	// Redirect
	http.Redirect(w, r, shortlink.Origin, http.StatusFound)
}

func ShortlinkCreationHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("ShortlinkCreationHandler")

	//Load params
	vars := mux.Vars(r)
	origin := vars["value"]
	token := r.FormValue("custom")
	origin = "http://" + origin

	token, err := RegisterShortlink(origin, token)
	if err != nil {
		http.Error(w, "Invalid origin parameter", 404)
		return
	}

	w.Write([]byte(origin + " is now accessible at the url http://localhost:8000/" + token + "\n"))
}

func AdvancedShortlinkCreationHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("AdvancedShortlinkCreationHandler")

	//Load params
	decoder := json.NewDecoder(r.Body)
	var request AdvancedShortlinkRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	origin := request.Origin
	token := request.Token

	token, err = RegisterShortlink(origin, token)
	if err != nil {
		http.Error(w, "Invalid origin parameter", 404)
		return
	}

	w.Write([]byte(origin + " is now accessible at the url http://localhost:8000/" + token + "\n"))
}

func RegisterShortlink(origin, token string) (string, error) {
	//Check origin
	_, err := http.Get(origin)
	if err != nil {
		return "", err
	}

	//Store in Redis
	if token == "" {
		token = RandomString()
	}

	//check if key is already used
	if redisClient.Get(token).Val() != "" {
		log.WithFields(log.Fields{"token": token}).Warn("token already used")
		i := 0
		for ; redisClient.Get(token+strconv.Itoa(i)).Val() != ""; i++ {
			log.WithFields(log.Fields{"token": token + strconv.Itoa(i)}).Warn("token already used")
		}
		token = token + strconv.Itoa(i)
	}

	// Save Shortlink in Redis
	log.WithFields(log.Fields{"origin": origin, "token": token}).Info("creation")
	uuid, _ := newUUID()
	CreateShortlink(Shortlink{uuid, token, origin, time.Now().Unix(), 0})
	return token, nil
}

func MonitoringHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("MonitoringHandler")
	//Load params
	vars := mux.Vars(r)
	token := vars["value"]
	shortlink, err := ReadFromRedisAsJson(token)
	if err != nil {
		http.Error(w, "Token not found", 404)
		return
	}
	w.Write([]byte(shortlink))
}

// UTILS

// Save a Shortlink object
func CreateShortlink(shortlink Shortlink) {
	// Shortlink -> JSON
	shortlinkAsJson, err := json.Marshal(shortlink)
	if err != nil {
		log.Fatal("could not Marshall a shortlink")
	}
	// Save Key
	redisClient.SetNX(shortlink.Token, string(shortlinkAsJson), 5*time.Minute)
}

// Update a Shortlink object
func UpdateShortlink(shortlink Shortlink) {
	// Shortlink -> JSON
	shortlinkAsJson, err := json.Marshal(shortlink)
	if err != nil {
		log.Fatal("could not Marshall a shortlink")
	}
	// Update Shortlink
	redisClient.SetXX(shortlink.Token, string(shortlinkAsJson), 5*time.Minute)
}

// Get a Shortlink object
func ReadFromRedis(token string) (Shortlink, error) {
	// Get value from key
	redisValue := redisClient.Get(token).Val()
	if redisValue == "" {
		return Shortlink{"", "", "", 0, 0}, errors.New("No value is associated to key [" + token + "]")
	}
	// JSON -> Shortlink
	reader := json.NewDecoder(strings.NewReader(redisValue))
	var shortlink Shortlink
	reader.Decode(&shortlink)
	// Return
	return shortlink, nil
}

// Get a Shortlink object as it is stored in Redis (JSON)
func ReadFromRedisAsJson(token string) (string, error) {
	// Get value from key
	redisValue := redisClient.Get(token).Val()
	if redisValue == "" {
		return "", errors.New("No value is associated to key [" + token + "]")
	}
	// JSON -> Shortlink
	return redisValue, nil
}

// Generates a random String
func RandomString() string {
	rb := make([]byte, 6)
	dictionary := "abcdefghijklmnopqrstuvwxyz"
	rand.Read(rb)
	for k, v := range rb {
		rb[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(rb)
}

// Generates a UUID (http://play.golang.org/p/4FkNSiUDMg)
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
