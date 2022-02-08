package controller

import (
	"github.com/kiali/demos/topology-generator/pkg/api"
)

func Run(generatorConfig api.Generator, mode string) error {
	if mode == api.MODE_SERVER || mode == api.MODE_DEFAULT {
		if err := RunServer(); err != nil {
			return err
		}
	} else if mode == api.MODE_LOCAL {
		if err := RunCLI(generatorConfig); err != nil {
			return err
		}
	}
	return nil
}
