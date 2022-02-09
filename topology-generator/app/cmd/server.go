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

	"github.com/kiali/demos/topology-generator/pkg/controller"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "The web server for UI of Topology Generator",
	Long:  `The web server for the UI of Topology Generator.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := controller.RunServer(serverPort); err != nil {
			log.Fatalf("Run Server Error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "The port of Web Server")

	serverCmd.PersistentFlags().StringVarP(&name, "name", "N", "mimik", "Name of the instance")
	serverCmd.PersistentFlags().StringVarP(&istioProxyRequestCPU, "proxycpu", "C", "50m", "IstioProxy Request CPU")
	serverCmd.PersistentFlags().StringVarP(&istioProxyRequestMemory, "proxymem", "M", "128Mi", "IstioProxy Request Memory")
	serverCmd.PersistentFlags().StringVar(&mimikRequestCPU, "rcpu", "25m", "Mimik Request CPU")
	serverCmd.PersistentFlags().StringVar(&mimikRequestMemory, "rmem", "64Mi", "Mimik Request Memory")
	serverCmd.PersistentFlags().StringVar(&mimikLimitCPU, "lcpu", "200m", "Mimik Limit CPU")
	serverCmd.PersistentFlags().StringVar(&mimikLimitMemory, "lmem", "256Mi", "Mimik Limit Memory")
	serverCmd.PersistentFlags().StringVarP(&image, "image", "i", "quay.io/leandroberetta/mimik", "Image tag name")
	serverCmd.PersistentFlags().StringVarP(&version, "version", "v", "v0.0.2", "Image version")
	serverCmd.PersistentFlags().StringVarP(&enableInjection, "enable-injection", "e", "true", "Enable injection or not")
	serverCmd.PersistentFlags().StringVarP(&injectionlabel, "injection-label", "j", "istio-injection:enabled", "Injection Label")
	serverCmd.PersistentFlags().IntVar(&replicas, "replica", 1, "Number of Replicas created")
}
