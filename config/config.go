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

func InitConfig() (*Config, error) {
	config := &Config{}

	configFile := os.Getenv(EnvConfigFile)
	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read config file %q (set via %s): %w", configFile, EnvConfigFile, err)
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return nil, fmt.Errorf("invalid YAML in config file %q: %w", configFile, err)
		}
	} else {
		data, err := os.ReadFile(ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("config file not found at %s\n\nCreate it with:\n  mkdir -p ~/.config/netbird-cli\n  echo 'auth:\n  url: \"https://your-netbird-instance.example.com\"\n  token: \"your-api-token\"' > %s\n\nOr set environment variables:\n  export %s=<url>\n  export %s=<token>", ConfigFile, ConfigFile, EnvURL, EnvToken)
		}
		err = yaml.Unmarshal(data, config)
		if err != nil {
			return nil, fmt.Errorf("invalid YAML in config file %q: %w", ConfigFile, err)
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

	return config, nil
}
