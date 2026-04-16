package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Car struct {
	CarModel string  `json:"carModel"`
	Price    float32 `json:"price"`
}

type TravelInfo struct {
	City string `json:"city"`
	Cars []Car  `json:"cars"`
}

type Discount struct {
	User     string  `json:"user"`
	Discount float32 `json:"discount"`
}

var (
	currentService   = "cars"
	currentVersion   = "no-version"
	instance         = currentService + "/" + currentVersion
	listenAddress    = ":8091"
	discountsService = "http://localhost:8092"

	mysqlService  = "localhost:3306"
	mysqlUser     = "root"
	mysqlPassword = "password"
	mysqlDatabase = "test"

	chaosMonkey       = false
	chaosMonkeySleep  = 500 * time.Millisecond // Milliseconds to wait if chaosMonkey is enabled
	chaosMonkeyPortal = ""
	chaosMonkeyDevice = ""
	chaosMonkeyUser   = ""
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

func GetCars(w http.ResponseWriter, r *http.Request) {
	portal := r.Header.Get("portal")
	device := r.Header.Get("device")
	user := r.Header.Get("user")

	travelInfo, err, notFound := getTravelInfo(r)
	if err != nil {
		Error(w, notFound, err.Error())
		return
	}
	travelInfo = applyDiscounts(r, &travelInfo, "cars")

	glog.Infof("[%s] GetCars for city %s \n", instance, travelInfo.City)

	releaseTheMonkey(portal, device, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(travelInfo.Cars)
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

	results, err := db.Query("SELECT c2.carModel, c2.price FROM cities c1, cars c2 WHERE c1.cityId = c2.cityId AND c1.city = '" + cityName + "'")
	if err != nil {
		return TravelInfo{}, err, false
	}

	travelInfo := TravelInfo{}
	travelInfo.City = cityName
	travelInfo.Cars = make([]Car, 0)
	for results.Next() {
		var car Car
		err := results.Scan(&car.CarModel, &car.Price)
		if err != nil {
			glog.Errorf("[%s] getTravelInfo can't parse a car row: %s \n", instance, err.Error())
			continue
		}
		travelInfo.Cars = append(travelInfo.Cars, car)
	}
	if len(travelInfo.Cars) == 0 {
		return TravelInfo{}, errors.New("City " + cityName + " not found"), true
	}
	return travelInfo, nil, false
}

func applyDiscounts(r *http.Request, travelInfo *TravelInfo, discountFrom string) TravelInfo {
	user := r.Header.Get("user")
	if user == "" {
		return *travelInfo
	}

	discount := float32(1)
	request, _ := http.NewRequest("GET", discountsService+"/discounts/"+user, nil)
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

	for i, car := range travelInfo.Cars {
		travelInfo.Cars[i].Price = car.Price * discount
	}
	return *travelInfo
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
	// Keep business headers used by this demo logic.
	for _, header := range []string{"portal", "device", "user", "travel"} {
		value := a.Header.Get(header)
		if value != "" {
			b.Header.Set(header, value)
		}
	}

	// Extract trace context from inbound request and inject it outbound.
	propagator := propagation.NewCompositeTextMapPropagator(
		b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader|b3.B3SingleHeader)),
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	ctx := propagator.Extract(a.Context(), propagation.HeaderCarrier(a.Header))
	propagator.Inject(ctx, propagation.HeaderCarrier(b.Header))
}

func main() {
	setup()
	glog.Infof("Starting %s \n", instance)
	router := mux.NewRouter()
	router.HandleFunc("/cars/{city}", GetCars).Methods("GET")
	glog.Fatal(http.ListenAndServe(listenAddress, router))
}
