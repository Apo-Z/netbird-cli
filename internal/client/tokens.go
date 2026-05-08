package client

import (
	"fmt"
)

type Token struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ExpirationDate string `json:"expiration_date"`
	CreatedBy      string `json:"created_by"`
	CreatedAt      string `json:"created_at"`
	LastUsed       string `json:"last_used"`
}

type CreateTokenRequest struct {
	Name      string `json:"name"`
	ExpiresIn int    `json:"expires_in"`
}

type CreateTokenResponse struct {
	PlainToken         string `json:"plain_token"`
	PersonalAccessToken *Token `json:"personal_access_token"`
}

func (c *Client) GetTokens(userID string) ([]Token, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/users/%s/tokens", userID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetTokens : %w", err)
	}
	return bodyToSlice[Token](resp.Body)
}

func (c *Client) GetToken(userID, tokenID string) (*Token, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/users/%s/tokens/%s", userID, tokenID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetToken : %w", err)
	}
	return bodyToStructure[Token](resp.Body)
}

func (c *Client) CreateToken(userID string, req *CreateTokenRequest) (*CreateTokenResponse, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/users/%s/tokens", userID), req)
	if err != nil {
		return nil, fmt.Errorf("error CreateToken : %w", err)
	}
	return bodyToStructure[CreateTokenResponse](resp.Body)
}

func (c *Client) DeleteToken(userID, tokenID string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/users/%s/tokens/%s", userID, tokenID))
	if err != nil {
		return fmt.Errorf("error DeleteToken : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
