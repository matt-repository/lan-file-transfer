package config

type Config struct {
	ServerPort int
	DataDir    string
}

var config *Config

func Init(c *Config) {
	config = c
}

func Get() *Config {
	return config
}
