package client

import (
	"fmt"
)

type SetupKey struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	Expires              string   `json:"expires"`
	Type                 string   `json:"type"`
	Valid                bool     `json:"valid"`
	Revoked              bool     `json:"revoked"`
	UsedTimes            int      `json:"used_times"`
	LastUsed             string   `json:"last_used"`
	State                string   `json:"state"`
	AutoGroups           []string `json:"auto_groups"`
	UpdatedAt            string   `json:"updated_at"`
	UsageLimit           int      `json:"usage_limit"`
	Ephemeral            bool     `json:"ephemeral"`
	AllowExtraDNSLabels  bool     `json:"allow_extra_dns_labels"`
	Key                  string   `json:"key"`
}

type CreateSetupKeyRequest struct {
	Name                string   `json:"name"`
	Type                string   `json:"type"`
	ExpiresIn           int      `json:"expires_in"`
	AutoGroups          []string `json:"auto_groups"`
	UsageLimit          int      `json:"usage_limit"`
	Ephemeral           *bool    `json:"ephemeral,omitempty"`
	AllowExtraDNSLabels *bool    `json:"allow_extra_dns_labels,omitempty"`
}

type UpdateSetupKeyRequest struct {
	Revoked    bool     `json:"revoked"`
	AutoGroups []string `json:"auto_groups"`
}

func (c *Client) GetSetupKeys() ([]SetupKey, error) {
	resp, err := c.doGet("/api/setup-keys", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetSetupKeys : %w", err)
	}
	return bodyToSlice[SetupKey](resp.Body)
}

func (c *Client) GetSetupKeyByID(id string) (*SetupKey, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/setup-keys/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetSetupKeyByID : %w", err)
	}
	return bodyToStructure[SetupKey](resp.Body)
}

func (c *Client) GetSetupKeyByName(name string) (*SetupKey, error) {
	keys, err := c.GetSetupKeys()
	if err != nil {
		return nil, fmt.Errorf("error lookup setup-key %s : %w", name, err)
	}
	for _, k := range keys {
		if k.Name == name {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("setup key %s not found", name)
}

func (c *Client) CreateSetupKey(req *CreateSetupKeyRequest) (*SetupKey, error) {
	resp, err := c.doPost("/api/setup-keys", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateSetupKey : %w", err)
	}
	return bodyToStructure[SetupKey](resp.Body)
}

func (c *Client) UpdateSetupKey(id string, req *UpdateSetupKeyRequest) (*SetupKey, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/setup-keys/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateSetupKey : %w", err)
	}
	return bodyToStructure[SetupKey](resp.Body)
}

func (c *Client) DeleteSetupKey(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/setup-keys/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteSetupKey : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
