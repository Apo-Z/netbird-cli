package main

import (
	"fmt"
	"os"

	"netbird-cli/config"
	"netbird-cli/internal/client"
)

var c *client.Client

func main() {
	cfg := config.InitConfig()
	var err error
	c, err = client.NewNetbirdClient(cfg.Auth.URL, cfg.Auth.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API connection error: %s\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
