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
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project (i.e. namespace) to a cluster",
	Long: `Add a new project (i.e. namespace) to a cluster
with the required metadata for onboarding, i.e. domain owner,
team, support email, resource (cpu, memory) capacity.`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, _ := cmd.Flags().GetString("cluster")
		buildnumber, _ := cmd.Flags().GetString("buildnumber")
		name, _ := cmd.Flags().GetString("name")
		owner, _ := cmd.Flags().GetString("owner")
		team, _ := cmd.Flags().GetString("team")
		email, _ := cmd.Flags().GetString("email")
		cpu, _ := cmd.Flags().GetString("cpu")
		memory, _ := cmd.Flags().GetString("memory")

		if cluster == "" {
			fmt.Println("Please provide a valid cluster name")
			os.Exit(1)
		}
		fmt.Println("cluster:", cluster)
		fmt.Println("buildnumber:", buildnumber)
		fmt.Println("project name:", name)
		fmt.Println("owner:", owner)
		fmt.Println("team:", team)
		fmt.Println("email:", email)
		fmt.Println("cpu:", cpu)
		fmt.Println("memory:", memory)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP("cluster", "c", "", "cluster name (e.g. dnocp)")
	addCmd.Flags().StringP("buildnumber", "b", "", "build number (e.g. 003)")
	addCmd.Flags().StringP("name", "n", "", "project name (e.g. ingdirect-rc-test)")
	addCmd.Flags().StringP("owner", "o", "", "project owner")
	addCmd.Flags().StringP("team", "t", "", "team name")
	addCmd.Flags().StringP("email", "e", "", "support contact email")
	addCmd.Flags().StringP("cpu", "u", "2", "cpu requested capacity")
	addCmd.Flags().StringP("memory", "m", "16", "memory requested capacity")
}
