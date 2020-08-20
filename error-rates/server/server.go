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
)

var (
	listenAddress = ":8888"
	// HTTP Code,num requests
	codeRequests = "200,9;404,1"
	httpCodes = make([]int, 0)
	mapCodeRequests = make(map[int]int)
	counterRequests = make(map[int]int)
	serverName = "server"
	mutex sync.Mutex
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	la := os.Getenv("LISTEN_ADDRESS")
	if la != "" {
		listenAddress = la
		glog.Infof("LISTEN_ADDRESS=%s", listenAddress)
	}
	cr := os.Getenv("CODE_REQUESTS")
	if cr != "" {
		codeRequests = cr
		glog.Infof("CODE_REQUESTS=%s", codeRequests)
	}
	codes := strings.Split(codeRequests, ";")
	for _, code := range codes {
		split := strings.Split(code, ",")
		if len(split) == 2 {
			httpCode, _ := strconv.Atoi(split[0])
			counter, _ := strconv.Atoi(split[1])
			httpCodes = append(httpCodes, httpCode)
			mapCodeRequests[httpCode] = counter
			counterRequests[httpCode] = counter
		}
	}
}

func GetStatus(w http.ResponseWriter, _ *http.Request) {
	defer mutex.Unlock()
	mutex.Lock()
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
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(returnCode)
	status := "Status: " + strconv.Itoa(returnCode)
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
