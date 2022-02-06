package controller

import (
	"log"

	"github.com/leandroberetta/mimik/pkg/api"
)

func Run(generatorConfig api.Generator, mode string) {
	if mode == api.MODE_SERVER || mode == api.MODE_DEFAULT {
		if err := RunServer(); err != nil {
			log.Fatalf("Running Server error: %v", err)
		}
	} else if mode == api.MODE_LOCAL {
		if err := RunCLI(generatorConfig); err != nil {
			log.Fatalf("Running CLI error: %v", err)
		}
	}
}
