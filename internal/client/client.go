package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Client struct {
	*http.Client
	URL   string
	Token string

	mu        sync.RWMutex
	nameCache map[string]string
	idToName  map[string]string
}

func NewNetbirdClient(apiURL, token string) (*Client, error) {
	apiURL = strings.TrimRight(apiURL, "/")
	apiURL = strings.TrimSuffix(apiURL, "/api")

	if apiURL == "" {
		return nil, fmt.Errorf("API URL is not configured")
	}
	if token == "" {
		return nil, fmt.Errorf("API token is not configured")
	}

	client := &Client{
		Client:    &http.Client{},
		URL:       apiURL,
		Token:     token,
		nameCache: make(map[string]string),
		idToName:  make(map[string]string),
	}

	return client, nil
}

func (c *Client) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	reqURL := fmt.Sprintf("%s%s", c.URL, path)

	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("error request %s %s : %w", method, path, err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", c.Token))
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error HTTP %s %s : %w", method, path, err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		preview := string(body)
		if len(preview) > 300 {
			preview = preview[:300]
		}
		return nil, fmt.Errorf("%s %s -> HTTP %d : %s", method, path, resp.StatusCode, preview)
	}

	return resp, nil
}

func (c *Client) doGet(path string, queryParams map[string]string) (*http.Response, error) {
	if len(queryParams) > 0 {
		u, err := url.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("error parsing path : %w", err)
		}
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		path = u.String()
	}
	return c.doRequest("GET", path, nil)
}

func (c *Client) doPost(path string, payload interface{}) (*http.Response, error) {
	return c.doJSON("POST", path, payload)
}

func (c *Client) doPut(path string, payload interface{}) (*http.Response, error) {
	return c.doJSON("PUT", path, payload)
}

func (c *Client) doDelete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil)
}

func (c *Client) doJSON(method, path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshal payload : %w", err)
	}
	return c.doRequest(method, path, bytes.NewReader(data))
}

func bodyToStructure[T any](body io.ReadCloser) (*T, error) {
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var t T

	err = json.Unmarshal(data, &t)
	if err != nil {
		preview := string(data)
		if len(preview) > 300 {
			preview = preview[:300]
		}
		return nil, fmt.Errorf("invalid JSON response : %w\nbody=%s", err, preview)
	}

	return &t, nil
}

func bodyToSlice[E any](body io.ReadCloser) ([]E, error) {
	result, err := bodyToStructure[[]E](body)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func bodyToBytes(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return io.ReadAll(body)
}

func (c *Client) cacheLookup(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.nameCache[key]
	return v, ok
}

func (c *Client) cacheStore(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nameCache[key] = value
}

func (c *Client) GetRaw(path string) ([]byte, error) {
	resp, err := c.doGet(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) PutRaw(path string, body interface{}) ([]byte, error) {
	resp, err := c.doPut(path, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) PostRaw(path string, body interface{}) ([]byte, error) {
	resp, err := c.doPost(path, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) IDToName(resource, id string) string {
	c.mu.RLock()
	if name, ok := c.idToName[id]; ok {
		c.mu.RUnlock()
		return name
	}
	c.mu.RUnlock()

	var fetch func()
	switch resource {
	case "group":
		fetch = func() {
			groups, err := c.GetGroups()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, g := range groups {
				c.idToName[g.ID] = g.Name
			}
			c.mu.Unlock()
		}
	case "user":
		fetch = func() {
			users, err := c.GetUsers()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, u := range users {
				c.idToName[u.ID] = u.Name
				if u.Name == "" {
					c.idToName[u.ID] = u.Email
				}
			}
			c.mu.Unlock()
		}
	case "peer":
		fetch = func() {
			peers, err := c.GetPeers("", "")
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, p := range peers {
				c.idToName[p.ID] = p.Name
			}
			c.mu.Unlock()
		}
	case "network":
		fetch = func() {
			networks, err := c.GetNetworks()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, n := range networks {
				c.idToName[n.ID] = n.Name
			}
			c.mu.Unlock()
		}
	case "policy":
		fetch = func() {
			policies, err := c.GetPolicies()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, p := range policies {
				c.idToName[p.ID] = p.Name
			}
			c.mu.Unlock()
		}
	case "posturecheck":
		fetch = func() {
			checks, err := c.GetPostureChecks()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, pc := range checks {
				c.idToName[pc.ID] = pc.Name
			}
			c.mu.Unlock()
		}
	case "setupkey":
		fetch = func() {
			keys, err := c.GetSetupKeys()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, k := range keys {
				c.idToName[fmt.Sprintf("%d", k.ID)] = k.Name
			}
			c.mu.Unlock()
		}
	case "nameserver":
		fetch = func() {
			nsgs, err := c.GetNameserverGroups()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, nsg := range nsgs {
				c.idToName[nsg.ID] = nsg.Name
			}
			c.mu.Unlock()
		}
	case "dnszone":
		fetch = func() {
			zones, err := c.GetDNSZones()
			if err != nil {
				return
			}
			c.mu.Lock()
			for _, z := range zones {
				c.idToName[z.ID] = z.Name
			}
			c.mu.Unlock()
		}
	case "agentprovider":
		fetch = func() { c.ResolveAgentProviderID("") }
	case "agentpolicy":
		fetch = func() { c.ResolveAgentPolicyID("") }
	case "guardrail":
		fetch = func() { c.ResolveAgentGuardrailID("") }
	case "budgetrule":
		fetch = func() { c.ResolveAgentBudgetRuleID("") }
	default:
		return ""
	}

	fetch()

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.idToName[id]
}
