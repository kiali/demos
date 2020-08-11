package main

import (
	"net/http"
	"github.com/golang/glog"
	"flag"
	"os"
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
	"strings"
	"sync"
	"bytes"
	"sort"
)

const (

)

var (
	listenAddress = ":8080"
	// By default travel_control binary expects a /html folder in its directory
	webroot = "./html"
	statusDemoIndex = 0
	portalServices = make(map[string]string)
	portalStatus = make(map[string]PortalStatus)
	rw sync.Mutex

	errorCoordinates = []float64{-25.702539, 37.747909}
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	la := os.Getenv("LISTEN_ADDRESS")
	if la != "" {
		listenAddress = la
		glog.Infof("LISTEN_ADDRESS=%s", listenAddress)
	}

	ps := os.Getenv("PORTAL_SERVICES")
	if ps != "" {
		portals := strings.Split(ps, ",")
		for _, portal := range portals {
			details := strings.Split(portal, ";")
			if len(details) == 2 {
				portalServices[details[0]] = details[1]
			}
		}
		glog.Infof("PORTAL_SERVICES=%s", ps)
	} else {
		glog.Errorf("PORTAL_SERVICES is empty !! Travel Control won't start")
		os.Exit(1)
	}
	wr := os.Getenv("WEBROOT")
	if wr != "" {
		webroot = wr
		glog.Infof("WEBROOT=%s", webroot)
	}
}

func response(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		response, _ = json.Marshal(ResponseError{Error: err.Error()})
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func GetStatus(w http.ResponseWriter, _ *http.Request) {
	var wg sync.WaitGroup
	wg.Add(len(portalServices))

	for portalName, portalAddress := range portalServices {
		go func(name, address string) {
			defer wg.Done()
			request, _ := http.NewRequest("GET", address + "/status", nil)
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				glog.Errorf("Failed to get status for portal [%s] [%s]. Error: [%s]", name, address, err.Error())
				portalStatus[name] = PortalStatus{
					Name: name,
					Coordinates: errorCoordinates,
					Country: "Atlantis",
					Settings: Settings{},
					Status: Status{
						Error: true,
					},
				}
			} else {
				status := PortalStatus{}
				json.NewDecoder(response.Body).Decode(&status)
				glog.Infof("Received status for portal [%s]. Total requests: %d Cities: %d",name, status.Status.Requests.Total, len(status.Status.Cities))
				rw.Lock()
				portalStatus[name] = status
				rw.Unlock()
			}
		}(portalName, portalAddress)
	}

	wg.Wait()

	// Return always an ordered response
	portals := make([]string, 0)
	for portal := range portalStatus {
		portals = append(portals, portal)
	}

	sort.Strings(portals)

	status := make([]PortalStatus, 0)
	for _, portal := range portals {
		status = append(status, portalStatus[portal])
	}

	response(w, http.StatusOK, status)
}

func PutSettings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	portal := params["portal"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, ResponseError{Error: "Error reading body", Detail: err.Error()})
		return
	}

	var settings Settings
	if err := json.Unmarshal(body, &settings); err != nil {
		response(w, http.StatusBadRequest, ResponseError{Error: "Error unmarshall settings", Detail: err.Error()})
		return
	}

	if address, exist := portalServices[portal]; !exist {
		response(w, http.StatusNotFound, ResponseError{Error: "Error portal " + portal + " not found", Detail: err.Error()})
		return
	} else {
		request, _ := http.NewRequest("PUT", address + "/settings", bytes.NewBuffer(body))
		client := &http.Client{}
		_, err := client.Do(request)
		if err != nil {
			glog.Errorf("Failed to send a setting for portal [%s] [%s]. Error: [%s]", portal, address, err)
		}
	}

	glog.Infof("Received PutSettings for [%s] [%v]", portal, settings)
	response(w, http.StatusOK, settings)
}

func main() {
	setup()
	glog.Infof("Starting Travel Control")

	router := mux.NewRouter()

	// Dynamic routes

	router.HandleFunc("/status", GetStatus).Methods("GET")
	router.HandleFunc("/settings/{portal}", PutSettings).Methods("PUT")

	// Static routes

	// Travel Control console is available under "/console" suffix
	router.Path("/console").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webroot+"/index.html")
	})

	// Static resources are served from "/" path
	// Then index.html doesn't need to adapt the paths
	// i.e.
	// 		<link rel="stylesheet" href="style.css" />
	//		<script src="script.js"></script>
	//		d3.json("countries-110m.json")
	//
	// Don't forget to adjust this if it changes
	// Note that this setup allows to serve static from content and also to run it from built-in server in Goland
	staticFileServer := http.FileServer(http.Dir(webroot))
	router.PathPrefix("/").Handler(staticFileServer)

	glog.Fatal(http.ListenAndServe(listenAddress, router))
}
