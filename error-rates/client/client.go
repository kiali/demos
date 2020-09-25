package main

import (
	"github.com/golang/glog"
	"flag"
	"os"
	"net/http"
	"time"
	"strings"
	"io/ioutil"
	"strconv"
)

var (
	clientName = "client"
	serverURL = []string{"http://localhost:8888/status"}
	wait = 500
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	cn := os.Getenv("CLIENT_NAME")
	if cn != "" {
		clientName = cn
		glog.Infof("CLIENT_NAME=%s", clientName)
	}
	su := os.Getenv("SERVER_URL")
	if su != "" {
		serverURL = strings.Split(su, ",")
		glog.Infof("SERVER_URL=%s", serverURL)
	}
	w := os.Getenv("WAIT")
	if w != "" {
		if waitTime, err := strconv.Atoi(w); err == nil {
			wait = waitTime
			glog.Infof("WAIT=%s", w)
		}
	}
}

func main() {
	setup()
	glog.Infof("Starting Client Name [%s]", clientName)
	for {
		for _, url := range serverURL {
			request, _ := http.NewRequest("GET", url, nil)
			client := &http.Client{}
			response, err := client.Do(request)
			if err == nil {
				responseContent := ""
				if bodyBytes, err := ioutil.ReadAll(response.Body); err == nil {
					responseContent = string(bodyBytes)
					response.Body.Close()
				}
				glog.Infof("[%s] Request to %s - Code [%d] - Response [%s]", clientName, serverURL, response.StatusCode, responseContent)
			} else {
				glog.Infof("[%s] ERROR %s", clientName, err.Error())
			}
			requestSleep := time.Duration(wait) * time.Millisecond
			time.Sleep(requestSleep)
		}
	}
}
