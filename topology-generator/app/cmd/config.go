/*
Copyright © 2022 Xunzhuo <mixdeers@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"

	"github.com/kiali/demos/topology-generator/pkg/api"
	"github.com/kiali/demos/topology-generator/pkg/controller"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate topology config with commands",
	Long:  `Generate topology config with commands with Options.`,
	Run: func(cmd *cobra.Command, args []string) {

		generatorConfig = api.Generator{
			Namespaces:        namespaces,
			Services:          services,
			Connections:       connections,
			RandomConnections: randomconnections,
		}

		api.GlobalConfig = api.NewConfigurations(name, istioProxyRequestCPU, istioProxyRequestMemory, mimikRequestCPU,
			mimikRequestMemory, mimikLimitCPU, mimikLimitMemory, version, image, enableInjection, replicas, injectionlabel)

		if err := controller.RunCLI(generatorConfig, path); err != nil {
			log.Fatalf("Run Command Error: %v", err)
		}
	},
}

func init() {
	path, _ = os.Getwd()
	defaultPath := path + "/deploy.json"

	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().StringVarP(&name, "name", "N", "mimik", "Name of the instance")
	configCmd.PersistentFlags().StringVarP(&istioProxyRequestCPU, "proxycpu", "C", "50m", "IstioProxy Request CPU")
	configCmd.PersistentFlags().StringVarP(&istioProxyRequestMemory, "proxymem", "M", "128Mi", "IstioProxy Request Memory")
	configCmd.PersistentFlags().StringVar(&mimikRequestCPU, "rcpu", "25m", "Mimik Request CPU")
	configCmd.PersistentFlags().StringVar(&mimikRequestMemory, "rmem", "64Mi", "Mimik Request Memory")
	configCmd.PersistentFlags().StringVar(&mimikLimitCPU, "lcpu", "200m", "Mimik Limit CPU")
	configCmd.PersistentFlags().StringVar(&mimikLimitMemory, "lmem", "256Mi", "Mimik Limit Memory")
	configCmd.PersistentFlags().StringVarP(&image, "image", "i", "quay.io/leandroberetta/topogen", "Image tag name")
	configCmd.PersistentFlags().StringVarP(&version, "version", "v", ReleaseVersion, "Image version")
	configCmd.PersistentFlags().StringVarP(&enableInjection, "enable-injection", "e", "true", "Enable injection or not")
	configCmd.PersistentFlags().StringVarP(&injectionlabel, "injection-label", "j", "istio-injection:enabled", "Injection Label")
	configCmd.PersistentFlags().IntVar(&replicas, "replica", 1, "Number of Replicas created")

	configCmd.Flags().IntVarP(&namespaces, "namespace", "n", 5, "Number of Namespaces created")
	configCmd.Flags().IntVarP(&services, "service", "s", 5, "Number of Services created")
	configCmd.Flags().IntVarP(&connections, "connection", "c", 5, "Number of Connections created")
	configCmd.Flags().IntVarP(&randomconnections, "random", "r", 5, "Number of RandomConnections created")
	configCmd.Flags().StringVarP(&path, "path", "p", defaultPath, "Generated Config Path")
}
