package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var countriesGetCmd = &cobra.Command{
	Use:               "countries",
	Aliases:           []string{"co"},
	Short:             "List available country codes",
	ValidArgsFunction: validArgsFunc(countryNames),
	Run: func(cmd *cobra.Command, args []string) {
		countries, err := c.GetCountries()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(countries)
	},
}

var citiesGetCmd = &cobra.Command{
	Use:               "cities <country-code>",
	Aliases:           []string{"ci"},
	Short:             "List cities for a country",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(countryNames),
	Run: func(cmd *cobra.Command, args []string) {
		cities, err := c.GetCitiesByCountry(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(cities)
	},
}

func init() {
	getCmd.AddCommand(countriesGetCmd, citiesGetCmd)
}
