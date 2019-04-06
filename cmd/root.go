package cmd

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/bah2830/hp-5120-exporter/pkg/hpswitch"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hp-5120-exporter",
	Short: "Prometheus exporter for hp 5120 switches",
	Run: func(cmd *cobra.Command, args []string) {
		hpSwitch, err := hpswitch.NewWithPassword("localhost", 22002, "admin", "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer hpSwitch.Close()

		details, err := hpSwitch.GetEnvironmentDetails()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		spew.Dump(details)
		fmt.Println("here'")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("password", "", "password for ssh")
}
