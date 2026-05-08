package client

import (
	"fmt"
)

type Policy struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Enabled             bool     `json:"enabled"`
	SourcePostureChecks []string `json:"source_posture_checks"`
	Rules               []Rule   `json:"rules"`
}

type Rule struct {
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Enabled            bool                `json:"enabled"`
	Action             string              `json:"action"`
	Bidirectional      bool                `json:"bidirectional"`
	Protocol           string              `json:"protocol"`
	Ports              []string            `json:"ports,omitempty"`
	PortRanges         []PortRange         `json:"port_ranges,omitempty"`
	AuthorizedGroups   map[string][]string `json:"authorized_groups,omitempty"`
	ID                 string              `json:"id,omitempty"`
	Sources            []Group             `json:"sources,omitempty"`
	SourceResource     *ResourceRef        `json:"sourceResource,omitempty"`
	Destinations       []Group             `json:"destinations,omitempty"`
	DestinationResource *ResourceRef       `json:"destinationResource,omitempty"`
}

type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type CreatePolicyRequest struct {
	Name                string   `json:"name"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled"`
	SourcePostureChecks []string `json:"source_posture_checks,omitempty"`
	Rules               []Rule   `json:"rules"`
}

type UpdatePolicyRequest struct {
	Name                string   `json:"name"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled"`
	SourcePostureChecks []string `json:"source_posture_checks,omitempty"`
	Rules               []Rule   `json:"rules"`
}

func (c *Client) GetPolicies() ([]Policy, error) {
	resp, err := c.doGet("/api/policies", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetPolicies : %w", err)
	}
	return bodyToSlice[Policy](resp.Body)
}

func (c *Client) GetPolicyByID(id string) (*Policy, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/policies/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetPolicyByID : %w", err)
	}
	return bodyToStructure[Policy](resp.Body)
}

func (c *Client) GetPolicyByName(name string) (*Policy, error) {
	policies, err := c.GetPolicies()
	if err != nil {
		return nil, fmt.Errorf("error lookup policy %s : %w", name, err)
	}
	for _, p := range policies {
		if p.Name == name {
			c.cacheStore("policy:"+name, p.ID)
			c.cacheStore("policy:"+p.ID, p.ID)
			return &p, nil
		}
	}
	return nil, fmt.Errorf("policy %s not found", name)
}

func (c *Client) ResolvePolicyID(nameOrID string) (string, error) {
	if id, ok := c.cacheLookup("policy:" + nameOrID); ok {
		return id, nil
	}

	policy, err := c.GetPolicyByName(nameOrID)
	if err == nil {
		c.cacheStore("policy:"+policy.Name, policy.ID)
		c.cacheStore("policy:"+policy.ID, policy.ID)
		return policy.ID, nil
	}

	policy, err = c.GetPolicyByID(nameOrID)
	if err == nil {
		c.cacheStore("policy:"+policy.Name, policy.ID)
		c.cacheStore("policy:"+policy.ID, policy.ID)
		return policy.ID, nil
	}

	return "", fmt.Errorf("policy %s not found", nameOrID)
}

func (c *Client) CreatePolicy(req *CreatePolicyRequest) (*Policy, error) {
	resp, err := c.doPost("/api/policies", req)
	if err != nil {
		return nil, fmt.Errorf("error CreatePolicy : %w", err)
	}
	return bodyToStructure[Policy](resp.Body)
}

func (c *Client) UpdatePolicy(id string, req *UpdatePolicyRequest) (*Policy, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/policies/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdatePolicy : %w", err)
	}
	return bodyToStructure[Policy](resp.Body)
}

func (c *Client) DeletePolicy(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/policies/%s", id))
	if err != nil {
		return fmt.Errorf("error DeletePolicy : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
