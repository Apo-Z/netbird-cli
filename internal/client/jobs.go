package client

import "fmt"

type Job struct {
	ID           string            `json:"id"`
	CreatedAt    string            `json:"created_at"`
	CompletedAt  *string           `json:"completed_at"`
	TriggeredBy  string            `json:"triggered_by"`
	Status       string            `json:"status"`
	FailedReason *string           `json:"failed_reason"`
	Workload     map[string]interface{} `json:"workload"`
}

func (c *Client) GetPeerJobs(peerID string) ([]Job, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/peers/%s/jobs", peerID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetPeerJobs : %w", err)
	}
	return bodyToSlice[Job](resp.Body)
}

func (c *Client) GetPeerJob(peerID, jobID string) (*Job, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/peers/%s/jobs/%s", peerID, jobID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetPeerJob : %w", err)
	}
	return bodyToStructure[Job](resp.Body)
}

func (c *Client) CreatePeerJob(peerID string, workload map[string]interface{}) (*Job, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/peers/%s/jobs", peerID), workload)
	if err != nil {
		return nil, fmt.Errorf("error CreatePeerJob : %w", err)
	}
	return bodyToStructure[Job](resp.Body)
}
