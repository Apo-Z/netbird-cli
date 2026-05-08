package client

import (
	"fmt"
)

type Network struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Routers           []string `json:"routers"`
	RoutingPeersCount int      `json:"routing_peers_count"`
	Resources         []string `json:"resources"`
	Policies          []string `json:"policies"`
}

type CreateNetworkRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateNetworkRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type NetworkResource struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Enabled     bool    `json:"enabled"`
	Groups      []Group `json:"groups"`
}

type CreateNetworkResourceRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Address     string   `json:"address"`
	Enabled     bool     `json:"enabled"`
	Groups      []string `json:"groups"`
}

type UpdateNetworkResourceRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Address     string   `json:"address"`
	Enabled     bool     `json:"enabled"`
	Groups      []string `json:"groups"`
}

type NetworkRouter struct {
	ID         string   `json:"id"`
	Peer       string   `json:"peer"`
	PeerGroups []string `json:"peer_groups"`
	Metric     int      `json:"metric"`
	Masquerade bool     `json:"masquerade"`
	Enabled    bool     `json:"enabled"`
}

type CreateNetworkRouterRequest struct {
	Peer       string   `json:"peer,omitempty"`
	PeerGroups []string `json:"peer_groups,omitempty"`
	Metric     int      `json:"metric"`
	Masquerade bool     `json:"masquerade"`
	Enabled    bool     `json:"enabled"`
}

type UpdateNetworkRouterRequest struct {
	Peer       string   `json:"peer,omitempty"`
	PeerGroups []string `json:"peer_groups,omitempty"`
	Metric     int      `json:"metric"`
	Masquerade bool     `json:"masquerade"`
	Enabled    bool     `json:"enabled"`
}

func (c *Client) GetNetworks() ([]Network, error) {
	resp, err := c.doGet("/api/networks", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNetworks : %w", err)
	}
	return bodyToSlice[Network](resp.Body)
}

func (c *Client) GetNetworkByID(id string) (*Network, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/networks/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNetworkByID : %w", err)
	}
	return bodyToStructure[Network](resp.Body)
}

func (c *Client) GetNetworkByName(name string) (*Network, error) {
	networks, err := c.GetNetworks()
	if err != nil {
		return nil, fmt.Errorf("error lookup network %s : %w", name, err)
	}
	for _, n := range networks {
		if n.Name == name {
			c.cacheStore("network:"+name, n.ID)
			c.cacheStore("network:"+n.ID, n.ID)
			return &n, nil
		}
	}
	return nil, fmt.Errorf("network %s not found", name)
}

func (c *Client) ResolveNetworkID(nameOrID string) (string, error) {
	if id, ok := c.cacheLookup("network:" + nameOrID); ok {
		return id, nil
	}

	network, err := c.GetNetworkByName(nameOrID)
	if err == nil {
		c.cacheStore("network:"+network.Name, network.ID)
		c.cacheStore("network:"+network.ID, network.ID)
		return network.ID, nil
	}

	network, err = c.GetNetworkByID(nameOrID)
	if err == nil {
		c.cacheStore("network:"+network.Name, network.ID)
		c.cacheStore("network:"+network.ID, network.ID)
		return network.ID, nil
	}

	return "", fmt.Errorf("network %s not found", nameOrID)
}

func (c *Client) CreateNetwork(req *CreateNetworkRequest) (*Network, error) {
	resp, err := c.doPost("/api/networks", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateNetwork : %w", err)
	}
	return bodyToStructure[Network](resp.Body)
}

func (c *Client) UpdateNetwork(id string, req *UpdateNetworkRequest) (*Network, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/networks/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateNetwork : %w", err)
	}
	return bodyToStructure[Network](resp.Body)
}

func (c *Client) DeleteNetwork(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/networks/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteNetwork : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetNetworkResources(networkID string) ([]NetworkResource, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/networks/%s/resources", networkID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNetworkResources : %w", err)
	}
	return bodyToSlice[NetworkResource](resp.Body)
}

func (c *Client) GetNetworkResource(networkID, resourceID string) (*NetworkResource, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/networks/%s/resources/%s", networkID, resourceID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNetworkResource : %w", err)
	}
	return bodyToStructure[NetworkResource](resp.Body)
}

func (c *Client) CreateNetworkResource(networkID string, req *CreateNetworkResourceRequest) (*NetworkResource, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/networks/%s/resources", networkID), req)
	if err != nil {
		return nil, fmt.Errorf("error CreateNetworkResource : %w", err)
	}
	return bodyToStructure[NetworkResource](resp.Body)
}

func (c *Client) UpdateNetworkResource(networkID, resourceID string, req *UpdateNetworkResourceRequest) (*NetworkResource, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/networks/%s/resources/%s", networkID, resourceID), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateNetworkResource : %w", err)
	}
	return bodyToStructure[NetworkResource](resp.Body)
}

func (c *Client) DeleteNetworkResource(networkID, resourceID string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/networks/%s/resources/%s", networkID, resourceID))
	if err != nil {
		return fmt.Errorf("error DeleteNetworkResource : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetNetworkRouters(networkID string) ([]NetworkRouter, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/networks/%s/routers", networkID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNetworkRouters : %w", err)
	}
	return bodyToSlice[NetworkRouter](resp.Body)
}

func (c *Client) GetAllNetworkRouters() ([]NetworkRouter, error) {
	resp, err := c.doGet("/api/networks/routers", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAllNetworkRouters : %w", err)
	}
	return bodyToSlice[NetworkRouter](resp.Body)
}

func (c *Client) GetNetworkRouter(networkID, routerID string) (*NetworkRouter, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/networks/%s/routers/%s", networkID, routerID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNetworkRouter : %w", err)
	}
	return bodyToStructure[NetworkRouter](resp.Body)
}

func (c *Client) CreateNetworkRouter(networkID string, req *CreateNetworkRouterRequest) (*NetworkRouter, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/networks/%s/routers", networkID), req)
	if err != nil {
		return nil, fmt.Errorf("error CreateNetworkRouter : %w", err)
	}
	return bodyToStructure[NetworkRouter](resp.Body)
}

func (c *Client) UpdateNetworkRouter(networkID, routerID string, req *UpdateNetworkRouterRequest) (*NetworkRouter, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/networks/%s/routers/%s", networkID, routerID), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateNetworkRouter : %w", err)
	}
	return bodyToStructure[NetworkRouter](resp.Body)
}

func (c *Client) DeleteNetworkRouter(networkID, routerID string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/networks/%s/routers/%s", networkID, routerID))
	if err != nil {
		return fmt.Errorf("error DeleteNetworkRouter : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
