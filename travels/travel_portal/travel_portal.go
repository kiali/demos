package main

import (
	"flag"
	"os"
	"github.com/golang/glog"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"strings"
	"strconv"
	"sync"
	"time"
	mathRand "math/rand"
)

const (
	MAX_REQUEST_WAIT = 15000		// milliseconds
	MIN_REQUEST_WAIT = 750			// milliseconds
)

var (
	listenAddress     = ":8081"
	portalCoordinates = []float64{-5.321187, 35.890134}
	portalCountry     = "no-country"
	portalName        = "no-name"

	rw sync.RWMutex

	// Init settings
	settings = Settings{
		50,
		Devices{
			50,
			50,
		},
		Users{
			50,
			50,
		},
		TravelType{
			34,
			33,
			33,
		},
	}

	status = Status{
		Requests: Requests{
			Total: 0,
			Devices: Devices{
				0,
				0,
			},
			Users: Users{
				0,
				0,
			},
			TravelType: TravelType{
				0,
				0,
				0,
			},
		},
	}

	// Temp counters used to calculate next request type
	nextDevices = Devices{
		50,
		50,
	}
	nextUsers = Users{
		50,
		50,
	}
	nextTravelType = TravelType{
		34,
		33,
		33,
	}

	cities = []City{}

	travelsAgencyService = "http://localhost:8090/travels"
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	la := os.Getenv("LISTEN_ADDRESS")
	if la != "" {
		listenAddress = la
		glog.Infof("LISTEN_ADDRESS=%s", listenAddress)
	} else {
		glog.Warningf("LISTEN_ADDRESS variable empty. Using default [%s]", listenAddress)
	}
	pc := os.Getenv("PORTAL_COORDINATES")
	if pc != "" {
		split := strings.Split(pc, ",")
		if len(split) == 2 {
			var lat, lon float64
			var err error
			if lat, err = strconv.ParseFloat(split[0], 64); err != nil {
				glog.Errorf("PORTAL_COORDINATES cannot be parsed. Error in lat [%s]", lat)
			}
			if lon, err = strconv.ParseFloat(split[1], 64); err != nil {
				glog.Errorf("PORTAL_COORDINATES cannot be parsed. Error in lon [%s]", lon)
			}
			if lat != 0 && lon != 0 {
				// We flip the lat, lon -> lon, lat as d3 uses this format
				// But normal way to express this is lattitude, longitud (i.e. google maps)
				portalCoordinates = []float64{lon, lat}
				glog.Infof("PORTAL_COORDINATES=%s", pc)
			}
		}
	} else {
		glog.Warningf("PORTAL_COORDINATES variable empty. Using default [%f, %f]", portalCoordinates[1], portalCoordinates[0])
	}
	pu := os.Getenv("PORTAL_COUNTRY")
	if pu != "" {
		portalCountry = pu
		glog.Infof("PORTAL_COUNTRY=%s", portalCountry)
	} else {
		glog.Warningf("PORTAL_COUNTRY variable empty. Using default [%s]", portalCountry)
	}
	pn := os.Getenv("PORTAL_NAME")
	if pn != "" {
		portalName = pn
		glog.Infof("PORTAL_NAME=%s", portalName)
	} else {
		glog.Warningf("PORTAL_NAME variable empty. Using default [%s]", portalName)
	}
	tas := os.Getenv("TRAVELS_AGENCY_SERVICE")
	if tas != "" {
		travelsAgencyService = tas
		glog.Infof("TRAVELS_AGENCY_SERVICE=%s", travelsAgencyService)
	} else {
		glog.Warningf("TRAVELS_AGENCY_SERVICE variable empty. Using default [%s]", travelsAgencyService)
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

func calculateDevice(device string) string {
	if device == "mobile" {
		if nextDevices.Mobile > 0 {
			status.Requests.Devices.Mobile = status.Requests.Devices.Mobile + 1
			nextDevices.Mobile = nextDevices.Mobile - 1
		} else if nextDevices.Web > 0 {
			device = "web"
			status.Requests.Devices.Web = status.Requests.Devices.Web + 1
			nextDevices.Web = nextDevices.Web - 1
		}
	} else {
		if nextDevices.Web > 0 {
			status.Requests.Devices.Web = status.Requests.Devices.Web + 1
			nextDevices.Web = nextDevices.Web - 1
		} else {
			device = "mobile"
			status.Requests.Devices.Mobile = status.Requests.Devices.Mobile + 1
			nextDevices.Mobile = nextDevices.Mobile - 1
		}
	}
	if nextDevices.Mobile == 0 && nextDevices.Web == 0 {
		nextDevices.Mobile = settings.Devices.Mobile
		nextDevices.Web = settings.Devices.Web
	}
	return device
}

func calculateUser(user string) string {
	if user == "registered" {
		if nextUsers.Registered > 0 {
			status.Requests.Users.Registered = status.Requests.Users.Registered + 1
			nextUsers.Registered = nextUsers.Registered - 1
		} else if nextUsers.New > 0 {
			user = "new"
			status.Requests.Users.New = status.Requests.Users.New + 1
			nextUsers.New = nextUsers.New - 1
		}
	} else {
		if nextUsers.New > 0 {
			status.Requests.Users.New = status.Requests.Users.New + 1
			nextUsers.New = nextUsers.New - 1
		} else {
			user = "registered"
			status.Requests.Users.Registered = status.Requests.Users.Registered + 1
			nextUsers.Registered = nextUsers.Registered - 1
		}
	}
	if nextUsers.Registered == 0 && nextUsers.New == 0 {
		nextUsers.Registered = settings.Users.Registered
		nextUsers.New = settings.Users.New
	}
	return user
}

func calculateTravelType(travel_type string) string {
	if travel_type == "t1" {
		if nextTravelType.T1 > 0 {
			status.Requests.TravelType.T1 = status.Requests.TravelType.T1 +1
			nextTravelType.T1 = nextTravelType.T1 - 1
		} else if nextTravelType.T2 > 0 {
			travel_type = "t2"
			status.Requests.TravelType.T2 = status.Requests.TravelType.T2 +1
			nextTravelType.T2 = nextTravelType.T2 - 1
		} else if nextTravelType.T3 > 0 {
			travel_type = "t3"
			status.Requests.TravelType.T3 = status.Requests.TravelType.T3 +1
			nextTravelType.T3 = nextTravelType.T3 - 1
		}
	} else if travel_type == "t2" {
		if nextTravelType.T2 > 0 {
			status.Requests.TravelType.T2 = status.Requests.TravelType.T2 +1
			nextTravelType.T2 = nextTravelType.T2 - 1
		} else if nextTravelType.T3 > 0 {
			travel_type = "t3"
			status.Requests.TravelType.T3 = status.Requests.TravelType.T3 +1
			nextTravelType.T3 = nextTravelType.T3 - 1
		} else if nextTravelType.T1 > 0 {
			travel_type = "t1"
			status.Requests.TravelType.T1 = status.Requests.TravelType.T1 +1
			nextTravelType.T1 = nextTravelType.T1 - 1
		}
	} else {
		if nextTravelType.T3 > 0 {
			status.Requests.TravelType.T3 = status.Requests.TravelType.T3 +1
			nextTravelType.T3 = nextTravelType.T3 - 1
		} else if nextTravelType.T1 > 0 {
			travel_type = "t1"
			status.Requests.TravelType.T1 = status.Requests.TravelType.T1 +1
			nextTravelType.T1 = nextTravelType.T1 - 1
		} else if nextTravelType.T2 > 0 {
			travel_type = "t2"
			status.Requests.TravelType.T2 = status.Requests.TravelType.T2 +1
			nextTravelType.T2 = nextTravelType.T2 - 1
		}
	}
	if nextTravelType.T1 == 0 && nextTravelType.T2 == 0 && nextTravelType.T3 == 0 {
		nextTravelType.T1 = settings.TravelType.T1
		nextTravelType.T2 = settings.TravelType.T2
		nextTravelType.T3 = settings.TravelType.T3
	}
	return travel_type
}

func calculateRequestType() (string, string, string) {
	var device, user, travel_type string

	rw.Lock()

	// Update total requests for this portal
	status.Requests.Total = status.Requests.Total + 1

	// Devices and User index
	duIndex := status.Requests.Total % 2
	tIndex := status.Requests.Total % 3
	if duIndex == 0 {
		device = "mobile"
		user = "registered"

	} else {
		device = "web"
		user = "new"
	}
	if tIndex == 0 {
		travel_type = "t1"
	} else if tIndex == 1 {
		travel_type = "t2"
	} else {
		travel_type = "t3"
	}

	device = calculateDevice(device)
	user = calculateUser(user)
	travel_type = calculateTravelType(travel_type)

	rw.Unlock()

	return device, user, travel_type
}

func updateCityRequest(city string, coordinates []float64, device string, user string, travel_type string) {
	defer rw.Unlock()
	rw.Lock()

	var cityRequest *CityRequests

	var iCity int
	for i, cr := range status.Cities {
		if cr.City == city {
			cityRequest = &cr
			iCity = i
			break
		}
	}

	if cityRequest == nil {
		coordinates = []float64{0, 0}
		for _, c := range cities {
			if city == c.City {
				// We flip the lat, lon -> lon, lat as d3 uses this format
				// But normal way to express this is lattitude, longitud (i.e. google maps)
				coordinates[0], _ = strconv.ParseFloat(c.Lng, 64)
				coordinates[1], _ = strconv.ParseFloat(c.Lat, 64)
			}
		}
		cityRequest = &CityRequests{
			city,
			coordinates,
			Requests{
				0,
				Devices{},
				Users{},
				TravelType{},
			},
		}
		status.Cities = append(status.Cities, *cityRequest)
		iCity = len(status.Cities) - 1
	}

	status.Cities[iCity].Requests.Total = status.Cities[iCity].Requests.Total + 1

	if device == "mobile" {
		status.Cities[iCity].Requests.Devices.Mobile = status.Cities[iCity].Requests.Devices.Mobile + 1
	} else if device == "web" {
		status.Cities[iCity].Requests.Devices.Web = status.Cities[iCity].Requests.Devices.Web + 1
	}

	if user == "new" {
		status.Cities[iCity].Requests.Users.New = status.Cities[iCity].Requests.Users.New + 1
	} else if user == "registered" {
		status.Cities[iCity].Requests.Users.Registered = status.Cities[iCity].Requests.Users.Registered + 1
	}

	if travel_type == "t1" {
		status.Cities[iCity].Requests.TravelType.T1 = status.Cities[iCity].Requests.TravelType.T1 + 1
	} else if travel_type == "t2" {
		status.Cities[iCity].Requests.TravelType.T2 = status.Cities[iCity].Requests.TravelType.T2 + 1
	} else if travel_type == "t3" {
		status.Cities[iCity].Requests.TravelType.T3 = status.Cities[iCity].Requests.TravelType.T3 + 1
	}
}

func GetStatus(w http.ResponseWriter, _ *http.Request) {
	defer rw.RUnlock()
	glog.Infof("Sent Status for [%s]\n", portalName)

	rw.RLock()
	portalStatus := PortalStatus{
		portalName,
		portalCoordinates,
		portalCountry,
		settings,
		status,
	}
	response(w, http.StatusOK, portalStatus)
}

func PutSettings(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, ResponseError{Error: "Error reading body", Detail: err.Error()})
		return
	}

	var newSettings Settings
	if err := json.Unmarshal(body, &newSettings); err != nil {
		response(w, http.StatusBadRequest, ResponseError{Error: "Error unmarshall settings", Detail: err.Error()})
		return
	}
	rw.Lock()
	settings = newSettings
	nextDevices = settings.Devices
	nextUsers = settings.Users
	nextTravelType = settings.TravelType
	rw.Unlock()

	glog.Infof("Received New Settings for [%s] [%v]", portalName, settings)

	response(w, http.StatusOK, settings)
}

func main() {
	setup()
	glog.Infof("Starting Travel Portal [%s]", portalName)

	router := mux.NewRouter()

	router.HandleFunc("/status", GetStatus).Methods("GET")
	router.HandleFunc("/settings", PutSettings).Methods("PUT")

	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

	// Main loop
	go func() {
		for {
			device, user, travel := calculateRequestType()

			request, _ := http.NewRequest("GET", travelsAgencyService + "/travels", nil)
			request.Header.Set("portal", portalName)
			request.Header.Set("device", device)
			request.Header.Set("user", user)
			request.Header.Set("travel", travel)

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				glog.Errorf("Error requesting Destinations from [%s] - Device [%s] User [%s] Travel Type [%s] - Error: %s", portalName, device, user, travel, err.Error())
			}

			rw.Lock()
			json.NewDecoder(response.Body).Decode(&cities)
			response.Body.Close()
			rw.Unlock()

			if len(cities) > 0 {
				i := r.Int31n((int32)(len(cities)))
				city := cities[i].City

				request, _ := http.NewRequest("GET", travelsAgencyService + "/travels/" + city, nil)
				request.Header.Set("portal", portalName)
				request.Header.Set("device", device)
				request.Header.Set("user", user)
				request.Header.Set("travel", travel)

				client := &http.Client{}
				response, err := client.Do(request)
				if err != nil {
					glog.Errorf("Error requesting Travel Quote from [%s] - Device [%s] User [%s] Travel Type [%s] City [%s] - Error: %s", portalName, device, user, travel, city, err.Error())
				}
				if response.StatusCode >= 400 {
					glog.Errorf("Error requesting Travel Quote from [%s] - Device [%s] User [%s] Travel Type [%s] City [%s] - Status Error: %s", portalName, device, user, travel, city, response.StatusCode)

				} else {
					travelQuote := TravelQuote{}
					json.NewDecoder(response.Body).Decode(&travelQuote)
					updateCityRequest(city, travelQuote.Coordinates, device, user, travel)
					glog.Infof("[%s] Quote received - Device [%s] User [%s] Travel Type [%s] City [%s]. Quote %v", portalName, device, user, travel, city, travelQuote)
				}
			} else {
				glog.Warningf("Receiving empty Destinations from [%s] - Device [%s] User [%s] Travel Type [%s]", portalName, device, user, travel)
			}

			requestSleep := time.Duration(int(MAX_REQUEST_WAIT - (MAX_REQUEST_WAIT * float32(float32(settings.RequestRatio)/float32(100))))) * time.Millisecond

			// Protection against very fast requests ratios
			if requestSleep < MIN_REQUEST_WAIT * time.Millisecond {
				requestSleep = MIN_REQUEST_WAIT * time.Millisecond
			}

			glog.Infof("Sleep [%2.3f] s", requestSleep.Seconds())
			time.Sleep(requestSleep)
		}
	}()

	glog.Fatal(http.ListenAndServe(listenAddress, router))
}

