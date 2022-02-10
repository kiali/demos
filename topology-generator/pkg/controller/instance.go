package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiali/demos/topology-generator/pkg/api"
	"github.com/kiali/demos/topology-generator/pkg/service"
)

func RunInstance() error {
	log.Println("------------------------")
	log.Println("Running in Instance mode")
	log.Println("------------------------")
	instance, _ := service.NewService(
		os.Getenv("MIMIK_SERVICE_NAME"),
		os.Getenv("MIMIK_SERVICE_PORT"),
		os.Getenv("MIMIK_ENDPOINTS_FILE"),
		service.GetVersion(os.Getenv("MIMIK_LABELS_FILE")))

	client := &http.Client{}

	r := mux.NewRouter()
	r.Path("/").HandlerFunc(service.EndpointHandler(instance, client))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("serving at :8080")

	go srv.ListenAndServe()

	tc := make(chan struct{})
	if tg := os.Getenv("MIMIK_TRAFFIC_GENERATOR"); tg != "" {
		go generateTraffic(&instance, client, tc)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	close(tc)

	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("shutdown complete")
	return nil
}

func generateTraffic(service *api.Service, client *http.Client, quit chan struct{}) {
	for {
		select {
		case <-quit:
			log.Println("stopping traffic generator")
			return
		default:
			for _, endpoint := range service.Endpoints {
				req, _ := http.NewRequest(endpoint.Method, fmt.Sprintf("http://localhost:%d", 8080), nil)
				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
				}
				defer resp.Body.Close()
			}

			time.Sleep(1 * time.Second)
		}
	}
}
