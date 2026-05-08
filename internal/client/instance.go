package client

import "fmt"

type InstanceStatus struct {
	SetupRequired bool `json:"setup_required"`
}

type InstanceVersion struct {
	ManagementCurrentVersion   string `json:"management_current_version"`
	DashboardAvailableVersion  string `json:"dashboard_available_version"`
	ManagementAvailableVersion string `json:"management_available_version"`
	ManagementUpdateAvailable  bool   `json:"management_update_available"`
}

type SetupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	CreatePAT   *bool  `json:"create_pat,omitempty"`
	PatExpireIn *int   `json:"pat_expire_in,omitempty"`
}

type SetupResponse struct {
	UserID              string `json:"user_id"`
	Email               string `json:"email"`
	PersonalAccessToken string `json:"personal_access_token,omitempty"`
}

func (c *Client) GetInstanceStatus() (*InstanceStatus, error) {
	resp, err := c.doGet("/api/instance", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetInstanceStatus : %w", err)
	}
	return bodyToStructure[InstanceStatus](resp.Body)
}

func (c *Client) GetInstanceVersion() (*InstanceVersion, error) {
	resp, err := c.doGet("/api/instance/version", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetInstanceVersion : %w", err)
	}
	return bodyToStructure[InstanceVersion](resp.Body)
}

func (c *Client) SetupInstance(req *SetupRequest) (*SetupResponse, error) {
	resp, err := c.doPost("/api/setup", req)
	if err != nil {
		return nil, fmt.Errorf("error SetupInstance : %w", err)
	}
	return bodyToStructure[SetupResponse](resp.Body)
}
