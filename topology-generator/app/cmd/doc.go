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

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate the latest commmands docs",
	Long:  `Run this command to generate the latest command docs for topology generator.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("🚀 Start to generate Docs of Topology Generator")

		os.RemoveAll(docPath)

		if err := os.MkdirAll(docPath, os.ModePerm); err != nil {
			log.Fatalf("Create docPath err: %v", err)
		}

		if err := doc.GenMarkdownTree(rootCmd, docPath); err != nil {
			log.Fatal(err)
		}

		log.Println("👌 Finish to generate Docs of Topology Generator")
		log.Println("👀 You can find the latest generated docs in doc/commands")
	},
}

func init() {
	rootCmd.AddCommand(docCmd)
}
