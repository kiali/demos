package controller

import (
	"encoding/json"
	"log"
	"os"

	"github.com/kiali/demos/topology-generator/pkg/api"
	"github.com/kiali/demos/topology-generator/pkg/resources"
)

func RunCLI(generatorConfig api.Generator, path string) error {
	log.Println("Running in Local mode")
	log.Printf("Generating config:\n %+v", generatorConfig)
	log.Printf("Generating deploy config:\n %+v", api.GlobalConfig)

	topology := resources.GenerateTopology(generatorConfig, api.GlobalConfig)

	topo, err := json.Marshal(topology)

	if err != nil {
		return err
	}

	log.Printf("Writing Json to %s", path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(topo)

	if err != nil {
		return err
	}
	return nil
}
