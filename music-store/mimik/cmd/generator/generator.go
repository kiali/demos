package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/leandroberetta/mimik/pkg/api"
	"github.com/leandroberetta/mimik/pkg/resources"
)

func main() {
	r := mux.NewRouter()
	r.Path("/generate").HandlerFunc(generateTopologyHandler)
	r.Path("/").Handler(http.FileServer(http.Dir("./ui/build/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/build/static"))))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("serving at :8080")
	go srv.ListenAndServe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("shutdown complete")
	os.Exit(0)
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

	topology := resources.CreateTopology(generator.Namespaces, generator.Services, generator.Connections, generator.RandomConnections)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topology)
}
