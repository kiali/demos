package main

import (
	"github.com/golang/glog"
	"flag"
	"os"
	"net/http"
	"time"
	"strings"
)

var (
	clientName = "client"
	serverURL = []string{"http://localhost:8888/status"}
	wait = 500
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	su := os.Getenv("SERVER_URL")
	if su != "" {
		serverURL = strings.Split(su, ",")
		glog.Infof("SERVER_URL=%s", serverURL)
	}
}

func main() {
	setup()
	glog.Infof("Starting Client Name [%s]", clientName)
	for {
		for _, url := range serverURL {
			request, _ := http.NewRequest("GET", url, nil)
			client := &http.Client{}
			response, _ := client.Do(request)
			glog.Infof("[%s] Request to %s - Code [%d]", clientName, serverURL, response.StatusCode)
			requestSleep := time.Duration(wait) * time.Millisecond
			time.Sleep(requestSleep)
		}
	}
}
