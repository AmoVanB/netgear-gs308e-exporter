package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ServiceMetrics struct {
	TransmittedBytes *prometheus.GaugeVec
	ReceivedBytes    *prometheus.GaugeVec
	CRCErrorPackets  *prometheus.GaugeVec
	PortStatus       *prometheus.GaugeVec
	PortSpeed        *prometheus.GaugeVec
}

var ServiceMetricsVar ServiceMetrics

// InitMetrics sets new metrics and register them in Prometheus
func InitMetrics() {
	ServiceMetricsVar.TransmittedBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "tx_bytes",
			Help:      "The total number of bytes transmitted by a port of a Netgear GS308E switch",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(ServiceMetricsVar.TransmittedBytes)

	ServiceMetricsVar.ReceivedBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "rx_bytes",
			Help:      "The total number of bytes received by a port of a Netgear GS308E switch",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(ServiceMetricsVar.ReceivedBytes)

	ServiceMetricsVar.CRCErrorPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "crc_error_packets",
			Help:      "The total number of CRC error packets on a port of a Netgear GS308E switch",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(ServiceMetricsVar.CRCErrorPackets)

	ServiceMetricsVar.PortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "port_status",
			Help:      "1 if the port is Up, 0 otherwise",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(ServiceMetricsVar.PortStatus)

	ServiceMetricsVar.PortSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netgear_gs308e_exporter",
			Name:      "port_speed_mbps",
			Help:      "Linked speed of a port of a Netgear GS308E switch",
		},
		[]string{"port", "switch"},
	)
	prometheus.MustRegister(ServiceMetricsVar.PortSpeed)
}
