package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

// Redis client
var redisClient *redis.Client

// Shortlink structure
type Shortlink struct {
	Id, Token, Origin string
	CreationTs        int64
	Count             uint8
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

	log.Info("Registering MonitoringHandler on /admin/{value}")
	r.HandleFunc("/admin/{value}", MonitoringHandler).
		Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8000", r)
}

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

func RedirectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["value"]
	// Get Key in Redis
	shortlinkAsJson := redisClient.Get(token).Val()
	if shortlinkAsJson == "" {
		w.Write([]byte("invalid token!\n"))
		http.Error(w, http.StatusText(404), 404)
		return
	}
	dec := json.NewDecoder(strings.NewReader(shortlinkAsJson))
	var shortlink Shortlink
	dec.Decode(&shortlink)
	log.Info(shortlink)
	// TODO Increment count
	http.Redirect(w, r, shortlink.Origin, http.StatusFound)
}

func ShortlinkCreationHandler(w http.ResponseWriter, r *http.Request) {

	//Load params
	vars := mux.Vars(r)
	origin := vars["value"]
	token := r.FormValue("custom")
	origin = "http://" + origin
	log.WithFields(log.Fields{"origin": origin, "token": token}).Info("creation")
	//Check origin
	_, err := http.Get(origin)
	if err != nil {
		w.Write([]byte("invalid origin!\n"))
		http.Error(w, http.StatusText(404), 404)
		return
	}
	log.Info("origin is valid")
	//Store in Redis
	if token == "" {
		token = RandomString()
	}
	//check if key is already used
	log.Info("check token")

	if redisClient.Get(token).Val() != "" {
		log.WithFields(log.Fields{"token": token}).Warn("token already used")
		i := 0
		for ; redisClient.Get(token+strconv.Itoa(i)).Val() != ""; i++ {
			log.WithFields(log.Fields{"token": token + strconv.Itoa(i)}).Warn("token already used")
		}
		token = token + strconv.Itoa(i)
	}

	// Save Shortlink in Redis
	uuid, _ := newUUID()
	shortlinkAsJson, err := json.Marshal(Shortlink{uuid, token, origin, time.Now().Unix(), 0})
	if err != nil {
		log.Fatal(err)
	}
	redisClient.Append(token, string(shortlinkAsJson))
	w.Write([]byte(string(shortlinkAsJson) + "\n"))
}

func MonitoringHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("MonitoringHandler!\n"))
}

func RandomString() string {
	rb := make([]byte, 6)
	dictionary := "abcdefghijklmnopqrstuvwxyz"
	rand.Read(rb)
	for k, v := range rb {
		rb[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(rb)
}

//From http://play.golang.org/p/4FkNSiUDMg
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
