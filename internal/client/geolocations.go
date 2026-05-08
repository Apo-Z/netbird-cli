package client

import (
	"fmt"
)

type GeoCity struct {
	GeonameID int    `json:"geoname_id"`
	CityName  string `json:"city_name"`
}

func (c *Client) GetCountries() ([]string, error) {
	resp, err := c.doGet("/api/locations/countries", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetCountries : %w", err)
	}
	return bodyToSlice[string](resp.Body)
}

func (c *Client) GetCitiesByCountry(country string) ([]GeoCity, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/locations/countries/%s/cities", country), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetCitiesByCountry : %w", err)
	}
	return bodyToSlice[GeoCity](resp.Body)
}
