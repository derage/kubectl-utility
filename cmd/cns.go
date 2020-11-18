/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

// cnsCmd represents the cns command
var cnsCmd = &cobra.Command{
	Use:   "cns",
	Short: "change namespaces",
	Long:  `change the default namespace in your kube config`,
	Run: func(cmd *cobra.Command, args []string) {
		var namespaces Namespaces
		kubeCmd := exec.Command("kubectl", "get", "ns", "-o", "json")
		kubeCmd.Stderr = os.Stderr
		kubeNamespaces, err := kubeCmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error running get namespace:%s\n", err)
		}
		err = json.Unmarshal(kubeNamespaces, &namespaces)
		if err != nil {
			log.Fatalf("Error unmarshalling namespaces:%s\n", err)
		}
		idx, err := fuzzyfinder.Find(
			namespaces.Items,
			func(i int) string {
				return namespaces.Items[i].Metadata.Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				return fmt.Sprintf("Namespace: %s\nStatus: %s\n",
					namespaces.Items[i].Metadata.Name,
					namespaces.Items[i].Status.Phase)
			}),
		)

		currentContextCmd := exec.Command("kubectl", "config", "current-context")
		currentContext, err := currentContextCmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error running get namespace:%s\n", err)
		}

		setNamespacecmd := exec.Command("kubectl", "config", "set", "current-context", string(currentContext), "--namespace", namespaces.Items[idx].Metadata.Name)
		setNamespacecmd.Stdout = os.Stdout
		setNamespacecmd.Stderr = os.Stderr
		err = setNamespacecmd.Run()
		if err != nil {
			log.Fatalf("setting context failed: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cnsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cnsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cnsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Namespaces struct {
	APIVersion string      `json:"apiVersion"`
	Items      []Namespace `json:"items"`
	Kind       string      `json:"kind"`
}

type Namespace struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Labels            struct {
			Name    string `json:"name"`
			Updated string `json:"updated"`
		} `json:"labels"`
		Name string `json:"name"`
	} `json:"metadata,omitempty"`
	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`
}
