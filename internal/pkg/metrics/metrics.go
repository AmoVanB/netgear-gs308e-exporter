package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type exporterMetrics struct {
	TransmittedBytes *prometheus.GaugeVec
	ReceivedBytes    *prometheus.GaugeVec
	CRCErrorPackets  *prometheus.GaugeVec
	PortStatus       *prometheus.GaugeVec
	PortSpeed        *prometheus.GaugeVec
}

// Var holds the metrics of the program
var Var exporterMetrics

// InitMetrics defines the metrics and registers them in Prometheus
func InitMetrics() {
	Var.TransmittedBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "tx_bytes",
			Help:      "The total number of bytes transmitted by a port",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(Var.TransmittedBytes)

	Var.ReceivedBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "rx_bytes",
			Help:      "The total number of bytes received by a port",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(Var.ReceivedBytes)

	Var.CRCErrorPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "crc_error_packets",
			Help:      "The total number of CRC error packets on a port",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(Var.CRCErrorPackets)

	Var.PortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "port_status",
			Help:      "1 if the port is Up, 0 otherwise",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(Var.PortStatus)

	Var.PortSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "port_speed_mbps",
			Help:      "Linked speed of a port",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(Var.PortSpeed)
}
