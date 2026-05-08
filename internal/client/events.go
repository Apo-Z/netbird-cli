package client

import (
	"fmt"
)

type AuditEvent struct {
	ID             string            `json:"id"`
	Timestamp      string            `json:"timestamp"`
	Activity       string            `json:"activity"`
	ActivityCode   string            `json:"activity_code"`
	InitiatorID    string            `json:"initiator_id"`
	InitiatorName  string            `json:"initiator_name"`
	InitiatorEmail string            `json:"initiator_email"`
	TargetID       string            `json:"target_id"`
	Meta           map[string]string `json:"meta"`
}

type TrafficEvent struct {
	FlowID      string           `json:"flow_id"`
	ReporterID  string           `json:"reporter_id"`
	Source      *TrafficEndpoint `json:"source"`
	Destination *TrafficEndpoint `json:"destination"`
	User        *TrafficUser     `json:"user"`
	Policy      *TrafficPolicy   `json:"policy"`
	Protocol    int              `json:"protocol"`
	Direction   string           `json:"direction"`
	RXBytes     int64            `json:"rx_bytes"`
	RXPackets   int64            `json:"rx_packets"`
	TXBytes     int64            `json:"tx_bytes"`
	TXPackets   int64            `json:"tx_packets"`
}

type TrafficEndpoint struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	GeoLocation *GeoLocation `json:"geo_location,omitempty"`
	OS          string      `json:"os"`
	Address     string      `json:"address"`
	DNSLabel    string      `json:"dns_label"`
}

type TrafficUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type TrafficPolicy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PaginatedTrafficEvents struct {
	Data         []TrafficEvent `json:"data"`
	Page         int            `json:"page"`
	PageSize     int            `json:"page_size"`
	TotalRecords int            `json:"total_records"`
	TotalPages   int            `json:"total_pages"`
}

type ProxyAccessLog struct {
	ID              string            `json:"id"`
	ServiceID       string            `json:"service_id"`
	Timestamp       string            `json:"timestamp"`
	Method          string            `json:"method"`
	Host            string            `json:"host"`
	Path            string            `json:"path"`
	DurationMS      int               `json:"duration_ms"`
	StatusCode      int               `json:"status_code"`
	SourceIP        string            `json:"source_ip"`
	Reason          string            `json:"reason"`
	UserID          string            `json:"user_id"`
	AuthMethodUsed  string            `json:"auth_method_used"`
	CountryCode     string            `json:"country_code"`
	CityName        string            `json:"city_name"`
	SubdivisionCode string            `json:"subdivision_code"`
	BytesUpload     int64             `json:"bytes_upload"`
	BytesDownload   int64             `json:"bytes_download"`
	Protocol        string            `json:"protocol"`
	Metadata        map[string]string `json:"metadata"`
}

type PaginatedProxyLogs struct {
	Data         []ProxyAccessLog `json:"data"`
	Page         int              `json:"page"`
	PageSize     int              `json:"page_size"`
	TotalRecords int              `json:"total_records"`
	TotalPages   int              `json:"total_pages"`
}

func (c *Client) GetAuditEvents() ([]AuditEvent, error) {
	resp, err := c.doGet("/api/events/audit", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAuditEvents : %w", err)
	}
	return bodyToSlice[AuditEvent](resp.Body)
}

func (c *Client) GetTrafficEvents(params map[string]string) (*PaginatedTrafficEvents, error) {
	resp, err := c.doGet("/api/events/network-traffic", params)
	if err != nil {
		return nil, fmt.Errorf("error GetTrafficEvents : %w", err)
	}
	return bodyToStructure[PaginatedTrafficEvents](resp.Body)
}

func (c *Client) GetProxyAccessLogs(params map[string]string) (*PaginatedProxyLogs, error) {
	resp, err := c.doGet("/api/events/proxy", params)
	if err != nil {
		return nil, fmt.Errorf("error GetProxyAccessLogs : %w", err)
	}
	return bodyToStructure[PaginatedProxyLogs](resp.Body)
}
