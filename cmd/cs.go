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
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"

	// backwards engineered the config struct based on testing in kubectl pkg
	// https://github.com/kubernetes/kubectl/blob/7b01e2757cc74b1145976726b05cc1108ad2911d/pkg/cmd/config/use_context_test.go

	"k8s.io/client-go/tools/clientcmd"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

// csCmd represents the cs command
var csCmd = &cobra.Command{
	Use:   "cs",
	Short: "Change context of default kubernetes cluster",
	Long:  `This is a fuzzy finder quick way to change default kubernetes cluster you your kube config`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeCmd := exec.Command("kubectl", "config", "view")
		kubeConfig, err := kubeCmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error running kubectl config view:%s\n", err)
		}
		str := string(kubeConfig)

		config, err := clientcmd.Load([]byte(str))
		contexts := reflect.ValueOf(config.Contexts).MapKeys()
		_ = contexts
		idx, err := fuzzyfinder.Find(
			contexts,
			func(i int) string {
				return contexts[i].Interface().(string)
			})
		if err != nil {
			log.Fatal(err)
		}
		setContextTo := contexts[idx].Interface().(string)
		fmt.Printf("Setting context to %v\n", setContextTo)

		setContext := exec.Command("kubectl", "config", "set", "current-context", setContextTo)
		setContext.Stdout = os.Stdout
		setContext.Stderr = os.Stderr
		err = setContext.Run()
		if err != nil {
			log.Fatalf("setting context failed: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(csCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
