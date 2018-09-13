package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pusher/pusher-http-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DeploySource struct {
	Source      string `json: "source"`
	Comment     string `json: "comment"`
	Author      string `json: "author"`
	AuthorName  string `json: "author_name"`
	AuthorEmail string `json: "author_email"`
	DeployedAt  string `json: "deployed_at"`
	Revision    string `json: "revision"`
	Repository  string `json: "repository"`
	Environment string `json: "environment"`
	Server      string `json: "server"`
}

var PusherAppId = os.Getenv("PUSHER_APP_ID")
var PusherKey = os.Getenv("PUSHER_KEY")
var PusherSecret = os.Getenv("PUSHER_SECRET")
var PusherCluster = os.Getenv("PUSHER_CLUSTER")

func main() {

	if PusherAppId == "" {
		log.Fatal("Env PUSHER_APP_ID is Empty")
	}

	if PusherKey == "" {
		log.Fatal("Env PUSHER_KEY is Empty")
	}

	if PusherSecret == "" {
		log.Fatal("Env PUSHER_SECRET is Empty")
	}

	if PusherCluster == "" {
		log.Fatal("Env PUSHER_CLUSTER is Empty")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/deploy", DeployHandler).Methods("POST")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting Web Server port: 8000")

	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HOME")
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	var deploy DeploySource
	err := json.Unmarshal(bodyBytes, &deploy)
	if err != nil {
		log.Println("[ERROR] ", err)
	}

	w.WriteHeader(http.StatusOK)

	// instantiate a client
	client := pusher.Client{
		AppId:   PusherAppId,
		Key:     PusherKey,
		Secret:  PusherSecret,
		Cluster: PusherCluster,
	}

	data := map[string]string{"commit": deploy.Revision}
	// trigger an event on a channel, along with a data payload
	client.Trigger("deploy-notifications", "deploy", data)
	log.Println("New deploy notification send: ", string(bodyBytes))

	fmt.Fprintf(w, "DeployHandled")
}
