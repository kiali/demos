package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type City struct {
	City string `json:"city"`
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Flight struct {
	Airline string `json:"airline"`
	Price float32 `json:"price"`
}

type Hotel struct {
	Hotel string `json:"hotel"`
	Price float32 `json:"price"`
}

type Car struct {
	CarModel string `json:"carModel"`
	Price float32 `json:"price"`
}

type Insurance struct {
	Company string `json:"company"`
	Price float32 `json:"price"`
}

type TravelQuote struct {
	City string `json:"city"`
	Coordinates []float64 `json:"coordinates"`
	CreatedAt string `json:"createdAt"`
	Status string `json:"status"`
	Flights []Flight `json:"flights"`
	Hotels []Hotel `json:"hotels"`
	Cars []Car `json:"cars"`
	Insurances []Insurance `json:"insurances"`
}

var (
	currentService = "travels"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":8090"

	carsService = "http://localhost:8091"
	flightsService = "http://localhost:8093"
	hotelsService = "http://localhost:8094"
	insurancesService = "http://localhost:8095"

	chaosMonkey = false
	chaosMonkeySleep = 500 * time.Millisecond // Milliseconds to wait if chaosMonkey is enabled
	chaosMonkeyCity = ""
	chaosMonkeyUser = ""
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	ss := os.Getenv("CURRENT_SERVICE")
	if ss != "" {
		currentService = ss
	}
	sv := os.Getenv("CURRENT_VERSION")
	if sv != "" {
		currentVersion = sv
	}
	instance = currentService + "/" + currentVersion

	la := os.Getenv("LISTEN_ADDRESS")
	if la != "" {
		listenAddress = la
		glog.Infof("LISTEN_ADDRESS=%s", listenAddress)
	} else {
		glog.Warningf("LISTEN_ADDRESS variable empty. Using default [%s]", listenAddress)
	}

	fs := os.Getenv("FLIGHTS_SERVICE")
	if fs != "" {
		flightsService = fs
	}
	hs := os.Getenv("HOTELS_SERVICE")
	if hs != "" {
		hotelsService = hs
	}
	cs := os.Getenv("CARS_SERVICE")
	if cs != "" {
		carsService = cs
	}
	is := os.Getenv("INSURANCES_SERVICE")
	if is != "" {
		insurancesService = is
	}

	if os.Getenv("CHAOS_MONKEY") == "true" {
		chaosMonkey = true
		sleep := os.Getenv("CHAOS_MONKEY_SLEEP")
		if value, err := strconv.Atoi(sleep); err == nil {
			chaosMonkeySleep = time.Duration(value) * time.Millisecond
		}
		chaosMonkeyCity = os.Getenv("CHAOS_MONKEY_CITY")
		chaosMonkeyUser = os.Getenv("CHAOS_MONKEY_USER")
	}
}

func Error(w http.ResponseWriter, notFound bool, msg string) {
	errorType := "Internal Error"
	if notFound {
		errorType = "NotFound"
	}
	glog.Infof("[%s] %s: %s \n", instance, errorType, msg)

	response, _ := json.Marshal(map[string]string{"error": msg})
	w.Header().Set("Content-Type", "application/json")

	errorCode := http.StatusInternalServerError
	if notFound {
		errorCode = http.StatusNotFound
	}
	w.WriteHeader(errorCode)
	w.Write(response)
}

func GetDestinations(w http.ResponseWriter, r *http.Request) {
	portal := r.Header.Get("portal")
	device := r.Header.Get("device")
	user := r.Header.Get("user")
	travel := r.Header.Get("travel")

	glog.Infof("[%s] GetDestinations from [%s]. Device [%s]. User [%s]. Travel [%s] \n", instance, portal, device, user, travel)

	request, _ := http.NewRequest("GET", hotelsService + "/hotels", nil)
	propagateHeaders(r, request)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		Error(w, false, "Error fetching destinations for portal [" + portal + "]")
		return
	}
	cities := make([]City, 0)
	json.NewDecoder(response.Body).Decode(&cities)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities)
}

func GetTravelQuote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	portal := r.Header.Get("portal")
	device := r.Header.Get("device")
	user := r.Header.Get("user")
	travel := r.Header.Get("travel")
	city := params["city"]

	glog.Infof("[%s] GetTravelQuote from [%s]. Device [%s]. User [%s]. Travel [%s] for [city: %s].\n", instance, portal, device, user, travel, city)

	travelQuote := TravelQuote{
		City: city,
		CreatedAt: time.Now().Format(time.RFC3339),
		Status: "Not valid",
	}

	if travel == "" {
		travel = "t1"
	}

	// Travel Type logic
	// T1:	Flight + Car + Hotel + Insurance
	// T2:	Flight + Hotel + Insurance
	// T3:	Car + Hotel + Insurance
	wg := sync.WaitGroup{}
	wg.Add(4)
	errChan := make(chan error, 4)

	go func() {
		defer wg.Done()
		if travel == "t1" || travel == "t2" {
			request, _ := http.NewRequest("GET", flightsService + "/flights/" + city, nil)
			propagateHeaders(r, request)
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				errChan <- err
				return
			}
			if response.StatusCode >= 400 {
				errChan <- errors.New(response.Status)
				return
			}
			json.NewDecoder(response.Body).Decode(&travelQuote.Flights)
		}
	}()

	go func() {
		defer wg.Done()
		request, _ := http.NewRequest("GET", hotelsService + "/hotels/" + city, nil)
		propagateHeaders(r, request)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			errChan <- err
			return
		}
		if response.StatusCode >= 400 {
			errChan <- errors.New(response.Status)
			return
		}
		json.NewDecoder(response.Body).Decode(&travelQuote.Hotels)
	}()

	go func() {
		defer wg.Done()
		if travel == "t1" || travel == "t3" {
			request, _ := http.NewRequest("GET", carsService + "/cars/" + city, nil)
			propagateHeaders(r, request)

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				errChan <- err
				return
			}
			if response.StatusCode >= 400 {
				errChan <- errors.New(response.Status)
				return
			}
			json.NewDecoder(response.Body).Decode(&travelQuote.Cars)
		}
	}()

	go func() {
		defer wg.Done()
		request, _ := http.NewRequest("GET", insurancesService + "/insurances/" + city, nil)
		propagateHeaders(r, request)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			errChan <- err
			return
		}
		if response.StatusCode >= 400 {
			errChan <- errors.New(response.Status)
			return
		}
		json.NewDecoder(response.Body).Decode(&travelQuote.Insurances)
	}()

	wg.Wait()
	if len(errChan) != 0 {
		Error(w, true, "Travel Quote for " + city + " not found")
		return
	}

	travelQuote.Status = "Valid"

	releaseTheMonkey(city, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelQuote)
}

func releaseTheMonkey(city, user string) {
	if chaosMonkey {
		glog.Infof("[%s] ChaosMonkey introduced %s \n", instance, chaosMonkeySleep.String())
		if (city != "" && city == chaosMonkeyCity) || (user != "" && user == chaosMonkeyUser) || (chaosMonkeyCity == "" && chaosMonkeyUser == "") {
			time.Sleep(chaosMonkeySleep)
		}
	}
}

func propagateHeaders(a *http.Request, b *http.Request) {
	headers := []string{
		"portal",
		"device",
		"user",
		"travel",
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-ot-span-context",
	}
	for _, header := range headers {
		b.Header.Add(header, a.Header.Get(header))
	}
}

func main() {
	setup()
	glog.Infof("Starting %s \n", instance)
	router := mux.NewRouter()
	router.HandleFunc("/travels", GetDestinations).Methods("GET")
	router.HandleFunc("/travels/{city}", GetTravelQuote).Methods("GET")
	glog.Fatal(http.ListenAndServe(listenAddress, router))
}


