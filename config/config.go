package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	EnvURL        = "NETBIRD_CLI_URL"
	EnvToken      = "NETBIRD_CLI_TOKEN"
	EnvConfigFile = "NETBIRD_CLI_CONFIG_FILE"
)

var ConfigFile string = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".config/netbird-cli/config.yaml")

type Auth struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type Config struct {
	Auth Auth `json:"auth"`
}

func InitConfig() *Config {
	config := &Config{}

	configFile := os.Getenv(EnvConfigFile)
	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			panic(err)
		}
	} else {
		data, err := os.ReadFile(ConfigFile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, config)
		if err != nil {
			panic(err)
		}
	}

	url := os.Getenv(EnvURL)
	if url != "" {
		config.Auth.URL = url
	}

	token := os.Getenv(EnvToken)
	if token != "" {
		config.Auth.Token = token
	}

	return config
}
