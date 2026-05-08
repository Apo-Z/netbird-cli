package client

import (
	"fmt"
)

type PostureCheck struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Checks      *PostureChecks `json:"checks,omitempty"`
}

type PostureChecks struct {
	NBVersionCheck        *NBVersionCheck        `json:"nb_version_check,omitempty"`
	OSVersionCheck        *OSVersionCheck        `json:"os_version_check,omitempty"`
	GeoLocationCheck      *GeoLocationCheck      `json:"geo_location_check,omitempty"`
	PeerNetworkRangeCheck *PeerNetworkRangeCheck `json:"peer_network_range_check,omitempty"`
	ProcessCheck          *ProcessCheck          `json:"process_check,omitempty"`
}

type NBVersionCheck struct {
	MinVersion string `json:"min_version"`
}

type OSVersionCheck struct {
	Android *MinVersion      `json:"android,omitempty"`
	IOS     *MinVersion      `json:"ios,omitempty"`
	Darwin  *MinVersion      `json:"darwin,omitempty"`
	Linux   *MinKernelVersion `json:"linux,omitempty"`
	Windows *MinKernelVersion `json:"windows,omitempty"`
}

type MinVersion struct {
	MinVersion string `json:"min_version"`
}

type MinKernelVersion struct {
	MinKernelVersion string `json:"min_kernel_version"`
}

type GeoLocationCheck struct {
	Locations []GeoLocation `json:"locations"`
	Action    string        `json:"action"`
}

type GeoLocation struct {
	CountryCode string `json:"country_code"`
	CityName    string `json:"city_name,omitempty"`
}

type PeerNetworkRangeCheck struct {
	Ranges []string `json:"ranges"`
	Action string   `json:"action"`
}

type ProcessCheck struct {
	Processes []ProcessPath `json:"processes"`
}

type ProcessPath struct {
	LinuxPath   string `json:"linux_path,omitempty"`
	MacPath     string `json:"mac_path,omitempty"`
	WindowsPath string `json:"windows_path,omitempty"`
}

type CreatePostureCheckRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Checks      *PostureChecks `json:"checks,omitempty"`
}

type UpdatePostureCheckRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Checks      *PostureChecks `json:"checks,omitempty"`
}

func (c *Client) GetPostureChecks() ([]PostureCheck, error) {
	resp, err := c.doGet("/api/posture-checks", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetPostureChecks : %w", err)
	}
	return bodyToSlice[PostureCheck](resp.Body)
}

func (c *Client) GetPostureCheckByID(id string) (*PostureCheck, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/posture-checks/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetPostureCheckByID : %w", err)
	}
	return bodyToStructure[PostureCheck](resp.Body)
}

func (c *Client) GetPostureCheckByName(name string) (*PostureCheck, error) {
	checks, err := c.GetPostureChecks()
	if err != nil {
		return nil, fmt.Errorf("error lookup posture-check %s : %w", name, err)
	}
	for _, pc := range checks {
		if pc.Name == name {
			return &pc, nil
		}
	}
	return nil, fmt.Errorf("posture check %s not found", name)
}

func (c *Client) CreatePostureCheck(req *CreatePostureCheckRequest) (*PostureCheck, error) {
	resp, err := c.doPost("/api/posture-checks", req)
	if err != nil {
		return nil, fmt.Errorf("error CreatePostureCheck : %w", err)
	}
	return bodyToStructure[PostureCheck](resp.Body)
}

func (c *Client) UpdatePostureCheck(id string, req *UpdatePostureCheckRequest) (*PostureCheck, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/posture-checks/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdatePostureCheck : %w", err)
	}
	return bodyToStructure[PostureCheck](resp.Body)
}

func (c *Client) DeletePostureCheck(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/posture-checks/%s", id))
	if err != nil {
		return fmt.Errorf("error DeletePostureCheck : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
