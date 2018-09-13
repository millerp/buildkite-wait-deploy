package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pusher/pusher-http-go"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {

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

	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HOME")
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	bodyString, _ := ioutil.ReadAll(r.Body)
	var deploy DeploySource
	err := json.Unmarshal(bodyString, &deploy)
	if err != nil {
		log.Println("[ERROR] ", err)
	}

	w.WriteHeader(http.StatusOK)

	// instantiate a client
	client := pusher.Client{
		AppId:   "387102",
		Key:     "349519f1474cd2dfcf8e",
		Secret:  "0d83d7603e89a5a855a0",
		Cluster: "us2",
	}

	data := map[string]string{"commit": deploy.Revision}
	// trigger an event on a channel, along with a data payload
	client.Trigger("deploy-notifications", "deploy", data)
	log.Println("New deploy notification send: ", deploy)

	fmt.Fprintf(w, "DeployHandled")
}
