package main

import (
	"encoding/json"
	"github.com/toorop/go-pusher"
	"log"
	"os"
	"time"
)

const (
	AppKey = "349519f1474cd2dfcf8e"
)

type DeployInfo struct {
	Commit string `json: "commit"`
}

var deploy DeployInfo

var commit = os.Getenv("BUILDKITE_COMMIT")
func main() {
	//flag.StringVar(&commit, "commit", os.Getenv("BUILDKITE_COMMIT"), "Commit ID")
	//flag.Parse()

INIT:
	log.Println("init...")
	pusherClient, err := pusher.NewCustomClient(AppKey, "ws-us2.pusher.com:443", "wss")
	if err != nil {
		log.Fatalln(err)
	}

	// Subscribe
	err = pusherClient.Subscribe("deploy-notifications")
	if err != nil {
		log.Println("Subscription error : ", err)
	}
	// Bind events
	deployEvents, err := pusherClient.Bind("deploy")
	if err != nil {
		log.Println("Bind error: ", err)
	}
	log.Println("Binded to 'deploy' event")

	// Test bind err
	errChannel, err := pusherClient.Bind(pusher.ErrEvent)
	if err != nil {
		log.Println("Bind error: ", err)
	}
	log.Println("Binded to 'ErrEvent' event")
	go LogWaitDeploy()
	for {
		select {
		case deployEvt := <-deployEvents:
			var stringBytes = []byte(deployEvt.Data)
			json.Unmarshal(stringBytes, &deploy)
			log.Println("INFO:", deploy.Commit)
			log.Println("INFO:", commit)

			if deploy.Commit == commit {
				log.Println("Deploy end.")
				os.Exit(0)
			}

			log.Println("Invalid commit")

		case errEvt := <-errChannel:
			log.Println("ErrEvent: " + errEvt.Data)
			pusherClient.Close()
			time.Sleep(time.Second)
			goto INIT
		}
	}
}

func LogWaitDeploy() {

	log.Println("Waiting deploy. Commit: ", commit)

	time.AfterFunc(5*time.Second, func() {
		LogWaitDeploy()
	})
}
