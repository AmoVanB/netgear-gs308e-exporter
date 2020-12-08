package config

// Config represents the configuration of the exporter
type Config struct {
	Switches []SwitchConfig
	Interval int
	Port     int
}

// SwitchConfig represents the configuration for a particular target switch
type SwitchConfig struct {
	URL      string
	Password string
}

// Var holds the running configuration of the program
var Var = Config{
	Switches: []SwitchConfig{},
	Interval: 30,
	Port:     8080,
}
