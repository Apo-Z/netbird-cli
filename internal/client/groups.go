package client

import (
	"fmt"
)

type Group struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	PeersCount     int                 `json:"peers_count"`
	ResourcesCount int                 `json:"resources_count"`
	Issued         string              `json:"issued"`
	Peers          []PeerRef           `json:"peers"`
	Resources      []ResourceRef       `json:"resources"`
}

type PeerRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResourceRef struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type CreateGroupRequest struct {
	Name      string        `json:"name"`
	Peers     []string      `json:"peers,omitempty"`
	Resources []ResourceRef `json:"resources,omitempty"`
}

type UpdateGroupRequest struct {
	Name      string        `json:"name"`
	Peers     []string      `json:"peers,omitempty"`
	Resources []ResourceRef `json:"resources,omitempty"`
}

func (c *Client) GetGroups(name ...string) ([]Group, error) {
	params := map[string]string{}
	if len(name) > 0 && name[0] != "" {
		params["name"] = name[0]
	}
	resp, err := c.doGet("/api/groups", params)
	if err != nil {
		return nil, fmt.Errorf("error GetGroups : %w", err)
	}
	return bodyToSlice[Group](resp.Body)
}

func (c *Client) GetGroupByID(id string) (*Group, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/groups/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetGroupByID : %w", err)
	}
	return bodyToStructure[Group](resp.Body)
}

func (c *Client) GetGroupByName(name string) (*Group, error) {
	groups, err := c.GetGroups(name)
	if err != nil {
		return nil, fmt.Errorf("error lookup group %s : %w", name, err)
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("group %s not found", name)
	}
	c.cacheStore("group:"+name, groups[0].ID)
	c.cacheStore("group:"+groups[0].ID, groups[0].ID)
	return &groups[0], nil
}

func (c *Client) ResolveGroupID(nameOrID string) (string, error) {
	if id, ok := c.cacheLookup("group:" + nameOrID); ok {
		return id, nil
	}

	group, err := c.GetGroupByName(nameOrID)
	if err == nil {
		c.cacheStore("group:"+group.Name, group.ID)
		c.cacheStore("group:"+group.ID, group.ID)
		return group.ID, nil
	}

	group, err = c.GetGroupByID(nameOrID)
	if err == nil {
		c.cacheStore("group:"+group.Name, group.ID)
		c.cacheStore("group:"+group.ID, group.ID)
		return group.ID, nil
	}

	return "", fmt.Errorf("group %s not found", nameOrID)
}

func (c *Client) CreateGroup(req *CreateGroupRequest) (*Group, error) {
	resp, err := c.doPost("/api/groups", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateGroup : %w", err)
	}
	group, err := bodyToStructure[Group](resp.Body)
	if err != nil {
		return nil, err
	}
	c.cacheStore("group:"+group.Name, group.ID)
	c.cacheStore("group:"+group.ID, group.ID)
	return group, nil
}

func (c *Client) UpdateGroup(id string, req *UpdateGroupRequest) (*Group, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/groups/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateGroup : %w", err)
	}
	group, err := bodyToStructure[Group](resp.Body)
	if err != nil {
		return nil, err
	}
	c.cacheStore("group:"+group.Name, group.ID)
	c.cacheStore("group:"+group.ID, group.ID)
	return group, nil
}

func (c *Client) DeleteGroup(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/groups/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteGroup : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
