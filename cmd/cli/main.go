package main

import (
	"fmt"
	"os"

	"netbird-cli/config"
	"netbird-cli/internal/client"
)

var c *client.Client

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %s\n", err)
		os.Exit(1)
	}

	if cfg.Auth.URL == "" {
		fmt.Fprintf(os.Stderr, "API URL is not configured.\n\n")
		fmt.Fprintf(os.Stderr, "Configure it via one of:\n")
		fmt.Fprintf(os.Stderr, "  • Config file %s:\n", config.ConfigFile)
		fmt.Fprintf(os.Stderr, "      auth:\n")
		fmt.Fprintf(os.Stderr, "        url: \"https://your-netbird-instance.example.com\"\n")
		fmt.Fprintf(os.Stderr, "        token: \"nbp_...\"\n")
		fmt.Fprintf(os.Stderr, "  • Environment variable: export %s=<url>\n", config.EnvURL)
		os.Exit(1)
	}

	if cfg.Auth.Token == "" {
		fmt.Fprintf(os.Stderr, "API token is not configured.\n\n")
		fmt.Fprintf(os.Stderr, "Configure it via one of:\n")
		fmt.Fprintf(os.Stderr, "  • Config file %s:\n", config.ConfigFile)
		fmt.Fprintf(os.Stderr, "      auth:\n")
		fmt.Fprintf(os.Stderr, "        url: \"https://your-netbird-instance.example.com\"\n")
		fmt.Fprintf(os.Stderr, "        token: \"nbp_...\"\n")
		fmt.Fprintf(os.Stderr, "  • Environment variable: export %s=<token>\n", config.EnvToken)
		os.Exit(1)
	}

	c, err = client.NewNetbirdClient(cfg.Auth.URL, cfg.Auth.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API connection error: %s\n", err)
		os.Exit(1)
	}

	edition, editionSource, err = resolveEdition(cfg.Edition, c.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %s\n", err)
		os.Exit(1)
	}
	applyEdition(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
