package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type City struct {
	City string `json:"city"`
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Hotel struct {
	Hotel string `json:"hotel"`
	Price float32 `json:"price"`
}

type TravelInfo struct {
	City string `json:"city"`
	Hotels []Hotel `json:"hotels"`
}

type Discount struct {
	User string `json:"user"`
	Discount float32 `json:"discount"`
}

var (
	currentService = "hotels"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":8094"
	discountsService = "http://localhost:8092"

	mysqlService = "localhost:3306"
	mysqlUser = "root"
	mysqlPassword = "password"
	mysqlDatabase = "test"

	chaosMonkey = false
	chaosMonkeySleep = 500 * time.Millisecond // Milliseconds to wait if chaosMonkey is enabled
	chaosMonkeyPortal = ""
	chaosMonkeyDevice = ""
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

	ds := os.Getenv("DISCOUNTS_SERVICE")
	if ds != "" {
		discountsService = ds
	}

	ms := os.Getenv("MYSQL_SERVICE")
	if ms != "" {
		mysqlService = ms
	}
	mu := os.Getenv("MYSQL_USER")
	if mu != "" {
		mysqlUser = mu
	}
	mp := os.Getenv("MYSQL_PASSWORD")
	if mp != "" {
		mysqlPassword = mp
	}
	md := os.Getenv("MYSQL_DATABASE")
	if md != "" {
		mysqlDatabase = md
	}

	if os.Getenv("CHAOS_MONKEY") == "true" {
		chaosMonkey = true
		sleep := os.Getenv("CHAOS_MONKEY_SLEEP")
		if value, err := strconv.Atoi(sleep); err == nil {
			chaosMonkeySleep = time.Duration(value) * time.Millisecond
		}
		chaosMonkeyPortal = os.Getenv("CHAOS_MONKEY_PORTAL")
		chaosMonkeyDevice = os.Getenv("CHAOS_MONKEY_DEVICE")
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

func GetCities(w http.ResponseWriter, _ *http.Request) {
	glog.Infof("[%s] GetCities \n", instance)

	dataSourceName := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlService + ")/" + mysqlDatabase
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		Error(w, false, err.Error())
		return
	}
	defer db.Close()

	cities := make([]City, 0)
	results, err := db.Query("SELECT city, lat, lng FROM cities")
	if err != nil {
		Error(w, false, err.Error())
		return
	}

	for results.Next() {
		var city City
		results.Scan(&city.City, &city.Lat, &city.Lng)
		cities = append(cities, city)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities)
}

func GetHotels(w http.ResponseWriter, r *http.Request) {
	portal := r.Header.Get("portal")
	device := r.Header.Get("device")
	user := r.Header.Get("user")
	travelInfo, err, notFound := getTravelInfo(r)
	if err != nil {
		Error(w, notFound, err.Error())
		return
	}
	travelInfo = applyDiscounts(r, &travelInfo, "hotels")

	glog.Infof("[%s] GetHotels for city %s \n", instance, travelInfo.City)

	releaseTheMonkey(portal, device, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelInfo.Hotels)
}

func applyDiscounts(r *http.Request, travelInfo *TravelInfo, discountFrom string) TravelInfo {
	user := r.Header.Get("user")

	if user == "" {
		return *travelInfo
	}

	discount := float32(1)
	request, _ := http.NewRequest("GET", discountsService + "/discounts/" + user, nil)
	propagateHeaders(r, request)
	request.Header.Set("discountFrom", discountFrom)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		glog.Errorf("No discount. Discount service is not available")
		return *travelInfo
	}
	discountQuote := Discount{}
	json.NewDecoder(response.Body).Decode(&discountQuote)
	discount = discount - discountQuote.Discount

	glog.Infof("[%s] Applying discount %f for %s \n", instance, discount, user)

	for i, hotel := range travelInfo.Hotels {
		travelInfo.Hotels[i].Price = hotel.Price * discount
	}
	return *travelInfo
}

func getTravelInfo(r *http.Request) (TravelInfo, error, bool) {
	params := mux.Vars(r)
	cityName := params["city"]

	dataSourceName := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlService + ")/" + mysqlDatabase
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return TravelInfo{}, err, false
	}
	defer db.Close()

	results, err := db.Query("SELECT h2.hotel, h2.price FROM cities c1, hotels h2 WHERE c1.cityId = h2.cityId AND c1.city = '" + cityName + "'")
	if err != nil {
		return TravelInfo{}, err, false
	}

	travelInfo := TravelInfo{}
	travelInfo.City = cityName
	travelInfo.Hotels = make([]Hotel, 0)
	for results.Next() {
		var hotel Hotel
		err := results.Scan(&hotel.Hotel, &hotel.Price)
		if err != nil {
			glog.Errorf("[%s] getTravelInfo can't parse an hotel row %s \n", err.Error())
			continue
		}
		travelInfo.Hotels = append(travelInfo.Hotels, hotel)
	}
	if len(travelInfo.Hotels) == 0 {
		return TravelInfo{}, errors.New("City " + cityName + " not found"), true
	}
	return travelInfo, nil, false
}

func releaseTheMonkey(portal, device, user string) {
	if chaosMonkey {
		glog.Infof("[%s] ChaosMonkey introduced %s \n", instance, chaosMonkeySleep.String())
		if (portal != "" && portal == chaosMonkeyPortal) ||
			(device != "" && device == chaosMonkeyDevice) ||
			(user != "" && user == chaosMonkeyUser) ||
			(portal == "" && device == "" && user == "") {
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
	router.HandleFunc("/hotels", GetCities).Methods("GET")
	router.HandleFunc("/hotels/{city}", GetHotels).Methods("GET")
	log.Fatal(http.ListenAndServe(listenAddress, router))
}
