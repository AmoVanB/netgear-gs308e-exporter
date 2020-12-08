package cmd

import (
	"fmt"
	"github.com/amovanb/netgear-gs308e-exporter/internal/app"
	"github.com/amovanb/netgear-gs308e-exporter/internal/pkg/config"
	"github.com/amovanb/netgear-gs308e-exporter/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	metrics.InitMetrics()
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start monitoring Netgear GS308E switches and exporting Prometheus metrics",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&config.Var)
		if err != nil {
			return err
		}

		for _, switchConfig := range config.Var.Switches {
			go func(c config.SwitchConfig) {
				err = app.MonitorSwitch(c, config.Var.Interval)
				if err != nil {
					log.Fatal(err)
				}
			}(switchConfig)
		}

		http.Handle("/metrics", promhttp.Handler())
		return http.ListenAndServe(fmt.Sprintf(":%d", config.Var.Port), nil)
	},
}
