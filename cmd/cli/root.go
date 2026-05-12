package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
	dryRun       bool
	editFlag     bool
)

var rootCmd = &cobra.Command{
	Use:     "netbird-cli",
	Short:   "CLI to interact with the Netbird API",
	Version: "v0.1.0",
}

var getCmd = &cobra.Command{
	Use:   "get <resource> [name]",
	Short: "List or display a resource",
}

var createCmd = &cobra.Command{
	Use:   "create <resource>",
	Short: "Create a resource",
}

var editCmd = &cobra.Command{
	Use:   "edit <resource> <name>",
	Short: "Edit a resource",
}

var deleteCmd = &cobra.Command{
	Use:   "delete <resource> <name>",
	Short: "Delete a resource",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (json, yaml)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without executing")

	createCmd.PersistentFlags().BoolVar(&editFlag, "edit", false, "Open editor to fill in the manifest")

	rootCmd.RegisterFlagCompletionFunc("output", staticCompletion([]string{"json", "yaml"}))

	rootCmd.AddCommand(getCmd, createCmd, editCmd, deleteCmd)
}

func exitErr(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s: %s\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	}
}

func dryRunCheck(payload interface{}) bool {
	if dryRun {
		fmt.Println("[dry-run] request that would be sent:")
		printJSON(payload)
		return true
	}
	return false
}

func dryRunMsg(msg string) bool {
	if dryRun {
		fmt.Printf("[dry-run] %s\n", msg)
		return true
	}
	return false
}
