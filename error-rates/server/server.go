package main

import (
	"github.com/golang/glog"
	"flag"
	"os"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"strings"
	"strconv"
	"time"
	"io/ioutil"
)

var (
	listenAddress = ":8888"
	// HTTP Code,num requests
	codeRequests = "200,9;404,1"
	httpCodes = make([]int, 0)
	mapCodeRequests = make(map[int]int)
	counterRequests = make(map[int]int)
	delayRequests = make(map[int]int)
	serverName = "server"
	mutex sync.Mutex
	serverURL = []string{}
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	la := os.Getenv("LISTEN_ADDRESS")
	if la != "" {
		listenAddress = la
		glog.Infof("LISTEN_ADDRESS=%s", listenAddress)
	}
	sn := os.Getenv("SERVER_NAME")
	if sn != "" {
		serverName = sn
		glog.Infof("SERVER_NAME=%s", serverName)
	}
	su := os.Getenv("SERVER_URL")
	if su != "" {
		serverURL = strings.Split(su, ",")
		glog.Infof("SERVER_URL=%s", serverURL)
	}
	cr := os.Getenv("CODE_REQUESTS")
	if cr != "" {
		codeRequests = cr
		glog.Infof("CODE_REQUESTS=%s", codeRequests)
	}
	codes := strings.Split(codeRequests, ";")
	for _, code := range codes {
		split := strings.Split(code, ",")
		lenSplit := len(split)
		if lenSplit >= 2 && lenSplit <= 3 {
			httpCode, _ := strconv.Atoi(split[0])
			counter, _ := strconv.Atoi(split[1])
			httpCodes = append(httpCodes, httpCode)
			mapCodeRequests[httpCode] = counter
			counterRequests[httpCode] = counter
			if lenSplit == 3 {
				delay, _ := strconv.Atoi(split[2])
				delayRequests[httpCode] = delay
			} else {
				delayRequests[httpCode] = 0
			}
		}
	}
}

func GetStatus(w http.ResponseWriter, _ *http.Request) {
	defer mutex.Unlock()
	mutex.Lock()
        // if configured, act as client of configured servers
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
			glog.Infof("[%s] Request to %s - Code [%d] - Response [%s]", serverName, serverURL, response.StatusCode, responseContent)
		} else {
			glog.Infof("[%s] ERROR %s", serverName, err.Error())
		}
	}
        // now send response
	returnCode := 200
	found := false
	for _, httpCode := range httpCodes {
		if counterRequests[httpCode] > 0 {
			glog.Infof("Counter Code [%d] = %d", httpCode, counterRequests[httpCode])
			returnCode = httpCode
			counterRequests[httpCode] -= 1
			found = true
			break
		}
	}
	// Reset counter
	if !found {
		for httpCode, counter := range mapCodeRequests {
			counterRequests[httpCode] = counter
		}
	}
	// Delay
	delay, ok := delayRequests[returnCode]
	if ok && delay > 0 {
		glog.Infof("[%s] Delay: [%d] ms", serverName, delay)
		requestSleep := time.Duration(delay) * time.Millisecond
		time.Sleep(requestSleep)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(returnCode)
	status := "Server: " + serverName + " Status: " + strconv.Itoa(returnCode)
	_, _ = w.Write([]byte(status))
	glog.Infof("[%s] - %s", serverName, status)
}

func main() {
	setup()
	glog.Infof("Starting Server Name [%s]", serverName)

	router := mux.NewRouter()
	router.HandleFunc("/status", GetStatus).Methods("GET")

	glog.Fatal(http.ListenAndServe(listenAddress, router))
}
