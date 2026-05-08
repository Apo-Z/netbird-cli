package client

import "fmt"

type IdentityProvider struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Issuer       string `json:"issuer"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type CreateIDPRequest struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Issuer       string `json:"issuer"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type UpdateIDPRequest struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Issuer       string `json:"issuer"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (c *Client) GetIdentityProviders() ([]IdentityProvider, error) {
	resp, err := c.doGet("/api/identity-providers", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetIdentityProviders : %w", err)
	}
	return bodyToSlice[IdentityProvider](resp.Body)
}

func (c *Client) GetIdentityProviderByID(id string) (*IdentityProvider, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/identity-providers/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetIdentityProviderByID : %w", err)
	}
	return bodyToStructure[IdentityProvider](resp.Body)
}

func (c *Client) GetIdentityProviderByName(name string) (*IdentityProvider, error) {
	idps, err := c.GetIdentityProviders()
	if err != nil {
		return nil, fmt.Errorf("error lookup idp %s : %w", name, err)
	}
	for _, idp := range idps {
		if idp.Name == name {
			return &idp, nil
		}
	}
	return nil, fmt.Errorf("identity provider %s not found", name)
}

func (c *Client) CreateIdentityProvider(req *CreateIDPRequest) (*IdentityProvider, error) {
	resp, err := c.doPost("/api/identity-providers", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateIdentityProvider : %w", err)
	}
	return bodyToStructure[IdentityProvider](resp.Body)
}

func (c *Client) UpdateIdentityProvider(id string, req *UpdateIDPRequest) (*IdentityProvider, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/identity-providers/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateIdentityProvider : %w", err)
	}
	return bodyToStructure[IdentityProvider](resp.Body)
}

func (c *Client) DeleteIdentityProvider(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/identity-providers/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteIdentityProvider : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
