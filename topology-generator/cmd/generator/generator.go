package main

import (
	"flag"
	"log"

	"github.com/kiali/demos/topology-generator/pkg/api"
	"github.com/kiali/demos/topology-generator/pkg/controller"
)

var name, istioProxyRequestCPU, istioProxyRequestMemory, mimikRequestCPU,
	mimikRequestMemory, mimikLimitCPU, mimikLimitMemory, image, version, injection, injectionlabel string

var replicas int

var mode string

var namespaces, services, connections, randomconnections int

var generatorConfig api.Generator

func init() {
	flag.StringVar(&mode, "m", api.MODE_SERVER, "Running Mode: l(local) or s(server)")

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

	api.GlobalConfig = api.NewConfigurations(name, istioProxyRequestCPU, istioProxyRequestMemory, mimikRequestCPU,
		mimikRequestMemory, mimikLimitCPU, mimikLimitMemory, version, image, injection, replicas, injectionlabel)

	generatorConfig = api.Generator{
		Namespaces:        namespaces,
		Services:          services,
		Connections:       connections,
		RandomConnections: randomconnections,
	}
}

func main() {
	flag.Parse()

	if err := controller.Run(generatorConfig, mode); err != nil {
		log.Fatalf("Running Error: %v", err)
	}
}
