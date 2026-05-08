package client

import (
	"fmt"
)

type Route struct {
	ID                  string   `json:"id"`
	NetworkType         string   `json:"network_type"`
	Description         string   `json:"description"`
	NetworkID           string   `json:"network_id"`
	Enabled             bool     `json:"enabled"`
	Peer                string   `json:"peer"`
	PeerGroups          []string `json:"peer_groups"`
	Network             string   `json:"network"`
	Domains             []string `json:"domains"`
	Metric              int      `json:"metric"`
	Masquerade          bool     `json:"masquerade"`
	Groups              []string `json:"groups"`
	KeepRoute           bool     `json:"keep_route"`
	AccessControlGroups []string `json:"access_control_groups"`
	SkipAutoApply       bool     `json:"skip_auto_apply"`
}

type CreateRouteRequest struct {
	Description         string   `json:"description"`
	NetworkID           string   `json:"network_id"`
	Enabled             bool     `json:"enabled"`
	Peer                string   `json:"peer,omitempty"`
	PeerGroups          []string `json:"peer_groups,omitempty"`
	Network             string   `json:"network,omitempty"`
	Domains             []string `json:"domains,omitempty"`
	Metric              int      `json:"metric"`
	Masquerade          bool     `json:"masquerade"`
	Groups              []string `json:"groups"`
	KeepRoute           bool     `json:"keep_route"`
	AccessControlGroups []string `json:"access_control_groups,omitempty"`
	SkipAutoApply       *bool    `json:"skip_auto_apply,omitempty"`
}

type UpdateRouteRequest struct {
	Description         string   `json:"description"`
	NetworkID           string   `json:"network_id"`
	Enabled             bool     `json:"enabled"`
	Peer                string   `json:"peer,omitempty"`
	PeerGroups          []string `json:"peer_groups,omitempty"`
	Network             string   `json:"network,omitempty"`
	Domains             []string `json:"domains,omitempty"`
	Metric              int      `json:"metric"`
	Masquerade          bool     `json:"masquerade"`
	Groups              []string `json:"groups"`
	KeepRoute           bool     `json:"keep_route"`
	AccessControlGroups []string `json:"access_control_groups,omitempty"`
	SkipAutoApply       *bool    `json:"skip_auto_apply,omitempty"`
}

func (c *Client) GetRoutes() ([]Route, error) {
	resp, err := c.doGet("/api/routes", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetRoutes : %w", err)
	}
	return bodyToSlice[Route](resp.Body)
}

func (c *Client) GetRouteByID(id string) (*Route, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/routes/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetRouteByID : %w", err)
	}
	return bodyToStructure[Route](resp.Body)
}

func (c *Client) GetRouteByDescription(description string) (*Route, error) {
	routes, err := c.GetRoutes()
	if err != nil {
		return nil, fmt.Errorf("error lookup route %s : %w", description, err)
	}
	for _, r := range routes {
		if r.Description == description {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("route %s not found", description)
}

func (c *Client) CreateRoute(req *CreateRouteRequest) (*Route, error) {
	resp, err := c.doPost("/api/routes", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateRoute : %w", err)
	}
	return bodyToStructure[Route](resp.Body)
}

func (c *Client) UpdateRoute(id string, req *UpdateRouteRequest) (*Route, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/routes/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateRoute : %w", err)
	}
	return bodyToStructure[Route](resp.Body)
}

func (c *Client) DeleteRoute(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/routes/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteRoute : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
