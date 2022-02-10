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
	"io/ioutil"
	"log"
	"os"

	"github.com/kiali/demos/topology-generator/pkg/kubectl"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the generated topology in the Cluster",
	Long:  `Run this command to clean the generated topology in the Cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Start to clean generated Topology, it may takes some time, please wait.\n")
		var err error

		if err = cleanNamespacesByFile(); err != nil {
			if err = cleanNamespacesBySelector(); err != nil {
				log.Fatalf("Clean generated mesh error: %v", err)
			}
		}
	},
}

// cleanNamespacesBySelector clean resources by selector
func cleanNamespacesBySelector() error {
	var out []byte
	var err error

	k := kubectl.New(&kubectl.Config{Bin: binary})

	if out, err = k.DeleteBySelector("generated-by=mimik"); err != nil {
		return err
	}

	log.Printf("Logging: \n +%v", string(out))
	return nil
}

// cleanNamespacesBySelector clean resources by file
func cleanNamespacesByFile() error {
	var out []byte

	k := kubectl.New(&kubectl.Config{Bin: binary})
	generatedFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer generatedFile.Close()

	bytes, _ := ioutil.ReadAll(generatedFile)
	if out, err = k.Delete(bytes); err != nil {
		return err
	}

	log.Printf("Logging: \n +%v", string(out))
	return nil
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVarP(&binary, "binary", "b", "kubectl", "The binary name to interact with Cluster API.")
}
