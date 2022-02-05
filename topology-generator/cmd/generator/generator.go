package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/leandroberetta/mimik/pkg/api"
	"github.com/leandroberetta/mimik/pkg/resources"
)

const MODE_SERVER = "s"

var name, istioProxyRequestCPU, istioProxyRequestMemory, mimikRequestCPU,
	mimikRequestMemory, mimikLimitCPU, mimikLimitMemory, image, version, injection, injectionlabel string

var replicas int

var mode, path string

var namespaces, services, connections, randomconnections int

func init() {
	flag.StringVar(&mode, "m", MODE_SERVER, "Running mode: l(local) or s(server)")

	flag.IntVar(&namespaces, "n", 5, "Number of Namespaces created")
	flag.IntVar(&services, "s", 5, "Number of Services created")
	flag.IntVar(&connections, "c", 5, "Number of Connections created")
	flag.IntVar(&randomconnections, "r", 5, "Number of RandomConnections created")

	flag.StringVar(&name, "name", "mimik", "")
	flag.StringVar(&istioProxyRequestCPU, "pcpu", "50m", "IstioProxy Request CPU")
	flag.StringVar(&istioProxyRequestMemory, "pmem", "128Mi", "IstioProxy Request Memory")
	flag.StringVar(&mimikRequestCPU, "rcpu", "25m", "Mimik Request CPU")
	flag.StringVar(&mimikRequestMemory, "rmem", "64Mi", "Mimik Request Memory")
	flag.StringVar(&mimikLimitCPU, "lcpu", "200m", "Mimik Limit CPU")
	flag.StringVar(&mimikLimitMemory, "lmem", "256Mi", "Mimik Limit Memory")
	flag.StringVar(&image, "image", "quay.io/leandroberetta/mimik", "Image tag name")
	flag.StringVar(&version, "version", "v0.0.2", "Image version")
	flag.StringVar(&injection, "injection", "true", "Enable injection or not")
	flag.StringVar(&injectionlabel, "injectionlabel", "istio-injection:enabled", "Injection Label")
	flag.IntVar(&replicas, "replica", 1, "Number of Replicas created")

	path, _ = os.Getwd()
}

func main() {
	flag.Parse()

	if mode == MODE_SERVER {
		log.Println("Running in Server mode")
		r := mux.NewRouter()
		r.Path("/generate").HandlerFunc(generateTopologyHandler)
		r.Path("/").Handler(http.FileServer(http.Dir("./ui/build/")))
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/build/static"))))

		srv := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}

		log.Println("serving at :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Starting Server error: %v", err)
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		log.Println("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}

		log.Println("shutdown complete")
		os.Exit(0)
	} else {
		log.Println("Running in Local mode")

		generator := api.Generator{
			Namespaces:        namespaces,
			Services:          services,
			Connections:       connections,
			RandomConnections: randomconnections,
		}

		log.Printf("Generating config:\n %+v", generator)

		config := api.NewConfigurations(name, istioProxyRequestCPU, istioProxyRequestMemory, mimikRequestCPU,
			mimikRequestMemory, mimikLimitCPU, mimikLimitMemory, version, image, injection, replicas, injectionlabel)

		log.Printf("Generating deploy config:\n %+v", config)

		topology := resources.GenerateTopology(generator, config)

		topo, err := json.Marshal(topology)
		if err != nil {
			log.Fatalf("Can not marshal struct to yaml: %v", err)
		}

		log.Printf("Writing yaml to %s", path+"/deploy.json\n")

		f, err := os.OpenFile(path+"/deploy.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		_, err = f.Write(topo)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func generateTopologyHandler(w http.ResponseWriter, r *http.Request) {
	generator := api.Generator{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&generator)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	config := api.NewConfigurations(name, istioProxyRequestCPU, istioProxyRequestMemory, mimikRequestCPU,
		mimikRequestMemory, mimikLimitCPU, mimikLimitMemory, version, image, injection, replicas, injectionlabel)

	topology := resources.GenerateTopology(generator, config)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(topology); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
