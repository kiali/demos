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

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "The instance running in topology",
	Long:  `Run this command to start the instance running in the Cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := controller.RunInstance(); err != nil {
			log.Fatalf("Starting Server error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(instanceCmd)
}
