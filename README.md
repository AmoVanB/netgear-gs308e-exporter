[![License](https://img.shields.io/github/license/grafana/grafana)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/amovanb/netgear-gs308e-exporter)](https://goreportcard.com/report/github.com/amovanb/netgear-gs308e-exporter)

## Prometheus exporter for the Netgear GS308E Switch

The [Netgear GS308E switch](https://www.netgear.com/support/product/gs308e.aspx) is a 8-port L2 Gigabit Ethernet switch. 
It has a web management interface, but does not support SNMP.

This exporter fetches port statistics and port stats from the web interface of one or more GS308E switch(es) and exposes them as Prometheus metrics.
 
 ### Usage
 
 #### Configuration file

The URL of the switches to monitor, their passwords, the frequency at which to monitor and the port on which to expose Prometheus metrics can be configured through a YAML config file.
See the [example config file](config/config.yaml) for documentation and default values.
 
 #### Docker
 
```shell script
docker pull amovanb/netgear-gs308e-exporter:latest
docker run docker.io/amovanb/netgear-gs308e-exporter:latest -c config.yaml
```

#### Locally

```shell script
git clone https://github.com/AmoVanB/netgear-gs308e-exporter.git
cd netgear-gs308e-exporter
go build -a -mod=vendor -o ./netgear-gs308e-exporter
./netgear-gs308e-exporter -c config.yaml
```

### Exported metrics

All metrics have two labels:
- `port`: port number (from 1 to 8)
- `switch`: host part of the URL of the switch

```
# HELP netgear_gs308e_exporter_crc_error_packets The total number of CRC error packets on a port
# TYPE netgear_gs308e_exporter_crc_error_packets gauge

# HELP netgear_gs308e_exporter_port_speed_mbps Linked speed of a port
# TYPE netgear_gs308e_exporter_port_speed_mbps gauge

# HELP netgear_gs308e_exporter_port_status 1 if the port is Up, 0 otherwise
# TYPE netgear_gs308e_exporter_port_status gauge

# HELP netgear_gs308e_exporter_rx_bytes The total number of bytes received by a port
# TYPE netgear_gs308e_exporter_rx_bytes gauge

# HELP netgear_gs308e_exporter_tx_bytes The total number of bytes transmitted by a port
# TYPE netgear_gs308e_exporter_tx_bytes gauge
```