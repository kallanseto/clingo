package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create project command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long: `Create a new project in the specified cluster with the specified parameters
`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, _ := cmd.Flags().GetString("cluster")
		buildnumber, _ := cmd.Flags().GetString("buildnumber")
		name, _ := cmd.Flags().GetString("name")
		owner, _ := cmd.Flags().GetString("owner")
		service, _ := cmd.Flags().GetString("service")
		application, _ := cmd.Flags().GetString("application")
		domain, _ := cmd.Flags().GetString("domain")
		team, _ := cmd.Flags().GetString("team")
		email, _ := cmd.Flags().GetString("email")
		namespacevip, _ := cmd.Flags().GetString("namespacevip")
		snatip, _ := cmd.Flags().GetString("snatip")
		cpu, _ := cmd.Flags().GetInt("cpu")
		memory, _ := cmd.Flags().GetInt("memory")

		if cluster == "" {
			fmt.Println("Please provide a valid cluster name")
			os.Exit(1)
		}
		if buildnumber == "" {
			fmt.Println("Please provide a valid cluster buildnumber")
			os.Exit(1)
		}

		p := &Project{
			Name:         strings.ToLower(name),
			Team:         team,
			Email:        email,
			Owner:        owner,
			Service:      service,
			Application:  application,
			Domain:       domain,
			Namespacevip: namespacevip,
			Snatip:       snatip,
			CPU:          cpu,
			Memory:       memory,
		}
		fmt.Println("Create project request for", p)

		basedir := cluster + "/" + buildnumber + "/"
		if err := p.newProjectDir(basedir); err != nil {
			log.Fatal(err)
		}
		p.allocateEgressIP(basedir + "/" + "egress-ip-allocations.json")
		if p.writeProjectManifests(basedir) {
			p.updateClusterKustomization(basedir)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	createCmd.PersistentFlags().String("cluster", "", "cluster name (e.g. dnocp)")
	createCmd.PersistentFlags().String("buildnumber", "", "build number (e.g. 003)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().String("name", "", "project name (e.g. ingdirect-rc-test)")
	createCmd.Flags().String("owner", "", "project owner")
	createCmd.Flags().String("service", "", "service name")
	createCmd.Flags().String("application", "", "application name")
	createCmd.Flags().String("domain", "", "domain (i.e.business function) name")
	createCmd.Flags().String("team", "", "team name")
	createCmd.Flags().String("email", "", "support contact email")
	createCmd.Flags().String("namespacevip", "", "VIP that has been allocated for the namespace")
	createCmd.Flags().String("snatip", "", "source NAT ip for the namespace VIP")
	createCmd.Flags().Int("cpu", 0, "cpu requested capacity")
	createCmd.Flags().Int("memory", 0, "memory requested capacity")
}
