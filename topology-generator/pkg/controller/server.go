package controller

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiali/demos/topology-generator/pkg/api"
	"github.com/kiali/demos/topology-generator/pkg/resources"
)

//go:embed *
var currentDir embed.FS

var buildDir, staticDir fs.FS

func RunServer(port int) error {

	log.Println("----------------------")
	log.Println("Running in Server mode")

	log.Println("----------------------")
	log.Println("Embed build dir")

	var err error
	buildDir, err = fs.Sub(currentDir, "build")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("----------------------")

	log.Println("Embed static dir")
	staticDir, err = fs.Sub(buildDir, "static")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("----------------------")

	// log.Println("Read Dir build: ")
	// dirEntries, _ := currentDir.ReadDir("build")
	// for _, de := range dirEntries {
	// 	log.Println(de.Name(), de.IsDir())
	// }
	// log.Println("---------")

	// log.Println("Read Dir build/static: ")
	// entries, _ := fs.ReadDir(buildDir, "static")
	// for _, de := range entries {
	// 	log.Println(de.Name(), de.IsDir())
	// }
	// log.Println("---------")

	// log.Println("Read Dir build/static/css: ")
	// entries, _ = fs.ReadDir(staticDir, "css")
	// for _, de := range entries {
	// 	log.Println(de.Name(), de.IsDir())
	// }
	// log.Println("---------")

	r := mux.NewRouter()
	r.Path("/generate").HandlerFunc(generateTopologyHandler)
	r.Path("/").Handler(http.FileServer(http.FS(buildDir)))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticDir))))
	serverPort := fmt.Sprintf(":%d", port)

	srv := &http.Server{
		Addr:    serverPort,
		Handler: r,
	}

	log.Println("serving at " + serverPort)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("shutdown complete")
	return nil
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

	topology := resources.GenerateTopology(generator, api.GlobalConfig)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(topology); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
