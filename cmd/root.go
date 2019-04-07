package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bah2830/switch-exporter/pkg/exporter"
	"github.com/bah2830/switch-exporter/pkg/metrics"

	"github.com/bah2830/switch-exporter/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "switch-exporter",
	Short: "Prometheus exporter for hp 5120 switches",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())
		if err := config.LoadConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		sig := registerSignals()
		c := config.GetConfig()
		go metrics.Serve(c)

		e := exporter.New(c)

		if err := e.Start(); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		defer e.Stop()

		<-sig
	},
}

func registerSignals() chan os.Signal {
	var sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT)
	return sig
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Duration("interval", 30*time.Second, "time betting polling checks")

	rootCmd.PersistentFlags().StringP("ssh.host", "H", "localhost", "host address for the switch")
	rootCmd.PersistentFlags().StringP("ssh.username", "u", "admin", "ssh username for the switch")
	rootCmd.PersistentFlags().Uint16P("ssh.port", "p", 22, "ssh port for the switch")
	rootCmd.PersistentFlags().StringP("ssh.password", "P", "", "password for ssh")

	rootCmd.PersistentFlags().Uint16("metrics.port", 9090, "Port for prometheus metrics")
	rootCmd.PersistentFlags().String("metrics.path", "/metrics", "Path prometheus metrics will be served on")
}
