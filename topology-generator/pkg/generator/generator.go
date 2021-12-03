package generator

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/leandroberetta/mimik/pkg/api"
)

// GenerateTopology generates a topology based on parameteres
func GenerateTopology(numServices, numConnections, numNamespaces, numRandomConnections int) map[string][]api.Service {
	m := make(map[string][]api.Service)

	for i := 1; i <= numNamespaces; i++ {
		n := generateNamespaceName(i)
		log.Printf("generating services for namespace: %s", n)
		m[n] = GenerateServices(numServices, numConnections)
	}

	// If there is one namespace, skip random connections
	if len(m) > 1 {
		GenerateRandomConnections(m, numRandomConnections)
	}

	return m
}

// GenerateRandomConnections connects services from different namespaces that do not contain connections to other services to avoid infinite loop calls
func GenerateRandomConnections(topology map[string][]api.Service, numRandomConnections int) {
	rand.Seed(time.Now().UnixNano())

	totalEmptyServices := 0
	emptyServicesByNamespace := make(map[string][]api.Service)
	for namespace, services := range topology {
		for _, service := range services {
			if len(service.Endpoints[0].Connections) == 0 {
				emptyServicesByNamespace[namespace] = append(emptyServicesByNamespace[namespace], service)
				totalEmptyServices++
			}
		}
	}

	if totalEmptyServices < numRandomConnections {
		numRandomConnections = totalEmptyServices
	}

	for i := 0; i < numRandomConnections; i++ {
		srcNs := getRandomNamespace(1, len(emptyServicesByNamespace))
		srcSvc := getRandomService(emptyServicesByNamespace[srcNs])

		dstNs := srcNs
		for dstNs == srcNs {
			dstNs = getRandomNamespace(1, len(emptyServicesByNamespace))
		}
		dstSvc := getRandomService(emptyServicesByNamespace[dstNs])

		randomConnection := api.Connection{
			Service: fmt.Sprintf("%s.%s", emptyServicesByNamespace[dstNs][dstSvc].Name, dstNs),
			Port:    8080,
			Path:    "/",
			Method:  "GET",
		}

		skip := false
		for _, connection := range emptyServicesByNamespace[dstNs][dstSvc].Endpoints[0].Connections {
			if connection.Service == fmt.Sprintf("%s.%s", emptyServicesByNamespace[srcNs][srcSvc].Name, srcNs) {
				log.Println("skipping connection because dst already has src as connection")
				skip = true
			}
		}
		for _, connection := range emptyServicesByNamespace[srcNs][srcSvc].Endpoints[0].Connections {
			if connection.Service == fmt.Sprintf("%s.%s", emptyServicesByNamespace[dstNs][dstSvc].Name, dstNs) {
				log.Println("skipping connection because src already has dst as connection")
				skip = true
			}
		}

		if !skip {
			log.Printf("adding random connection: %s.%s -> %s.%s", emptyServicesByNamespace[srcNs][srcSvc].Name, srcNs, emptyServicesByNamespace[dstNs][dstSvc].Name, dstNs)
			emptyServicesByNamespace[srcNs][srcSvc].Endpoints[0].Connections = append(emptyServicesByNamespace[srcNs][srcSvc].Endpoints[0].Connections, randomConnection)
		}
	}
}

// GenerateServices generates services for a namespace
func GenerateServices(numServices, numConnections int) []api.Service {
	last := 1
	ns := []api.Service{}
	for i := 1; i <= numServices; i++ {
		s := api.Service{
			Name:    fmt.Sprintf("a%d", i),
			Version: "v1",
			Endpoints: []api.Endpoint{
				{
					Path:        "/",
					Method:      "GET",
					Connections: []api.Connection{},
				},
			},
		}
		log.Printf("generating service: %s", s.Name)
		if last < numServices {
			for j := 1; j <= numConnections; j++ {
				c := api.Connection{
					Service: fmt.Sprintf("a%d", last+j),
					Port:    8080,
					Path:    "/",
					Method:  "GET",
				}
				log.Printf("adding connection: %s", c.Service)
				s.Endpoints[0].Connections = append(s.Endpoints[0].Connections, c)

			}
			last = last + numConnections
		}
		ns = append(ns, s)
	}

	return ns
}

func generateNamespaceName(numNamespace int) string {
	return fmt.Sprintf("n%d", numNamespace)
}

func getRandomNamespace(from, to int) string {
	numNamespace := from + rand.Intn(to)
	return generateNamespaceName(numNamespace)
}

func getRandomService(services []api.Service) int {
	return rand.Intn(len(services))
}
