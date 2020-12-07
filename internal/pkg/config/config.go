package config

type Config struct {
	Switches []SwitchConfig
	Interval int
	Port	 int
}

type SwitchConfig struct {
	Url      string
	Password string
}

var Var = Config{
	Switches: []SwitchConfig{},
	Interval: 30,
	Port: 8080,
}