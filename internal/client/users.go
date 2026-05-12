package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	ID              string       `json:"id"`
	Email           string       `json:"email"`
	Password        string       `json:"password"`
	Name            string       `json:"name"`
	Role            string       `json:"role"`
	Status          string       `json:"status"`
	LastLogin       string       `json:"last_login"`
	AutoGroups      []string     `json:"auto_groups"`
	IsCurrent       *bool        `json:"is_current"`
	IsServiceUser   *bool        `json:"is_service_user"`
	IsBlocked       *bool        `json:"is_blocked"`
	PendingApproval *bool        `json:"pending_approval"`
	Issued          string       `json:"issued"`
	IdpID           string       `json:"idp_id"`
	Permissions     *Permissions `json:"permissions"`
}

type Permissions struct {
	IsRestricted *bool                      `json:"is_restricted"`
	Modules      map[string]map[string]bool `json:"modules"`
}

type CreateUserRequest struct {
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	AutoGroups   []string `json:"auto_groups"`
	IsServiceUser bool    `json:"is_service_user"`
}

type UpdateUserRequest struct {
	Role       string   `json:"role"`
	AutoGroups []string `json:"auto_groups"`
	IsBlocked  *bool    `json:"is_blocked"`
}

type Invite struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	AutoGroups  []string `json:"auto_groups"`
	ExpiresAt   string   `json:"expires_at"`
	CreatedAt   string   `json:"created_at"`
	Expired     bool     `json:"expired"`
	InviteToken string   `json:"invite_token"`
}

type CreateInviteRequest struct {
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	Role       string   `json:"role"`
	AutoGroups []string `json:"auto_groups"`
	ExpiresIn  int      `json:"expires_in,omitempty"`
}

type RegenerateInviteRequest struct {
	ExpiresIn int `json:"expires_in,omitempty"`
}

type InviteInfo struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	ExpiresAt string `json:"expires_at"`
	Valid     bool   `json:"valid"`
	InvitedBy string `json:"invited_by"`
}

type AcceptInviteRequest struct {
	Password string `json:"password"`
}

type AcceptInviteResponse struct {
	Success bool `json:"success"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (c *Client) GetUsers(serviceUsersOnly ...bool) ([]User, error) {
	params := map[string]string{}
	if len(serviceUsersOnly) > 0 && serviceUsersOnly[0] {
		params["service_user"] = "true"
	}
	resp, err := c.doGet("/api/users", params)
	if err != nil {
		return nil, fmt.Errorf("error GetUsers : %w", err)
	}
	defer resp.Body.Close()

	var users []User
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Client) GetUserByEmail(email string) (*User, error) {
	users, err := c.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("error lookup user by email %s : %w", email, err)
	}
	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("user with email %s not found", email)
}

func (c *Client) CreateUser(req *CreateUserRequest) (*User, error) {
	resp, err := c.doPost("/api/users", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateUser : %w", err)
	}
	return bodyToStructure[User](resp.Body)
}

func (c *Client) UpdateUser(id string, req *UpdateUserRequest) (*User, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/users/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateUser : %w", err)
	}
	return bodyToStructure[User](resp.Body)
}

func (c *Client) DeleteUser(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/users/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteUser : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) ResendUserInvitation(id string) error {
	resp, err := c.doPost(fmt.Sprintf("/api/users/%s/invite", id), nil)
	if err != nil {
		return fmt.Errorf("error ResendUserInvitation : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) ApproveUser(id string) (*User, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/users/%s/approve", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error ApproveUser : %w", err)
	}
	return bodyToStructure[User](resp.Body)
}

func (c *Client) RejectUser(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/users/%s/reject", id))
	if err != nil {
		return fmt.Errorf("error RejectUser : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) ChangeUserPassword(id string, req *ChangePasswordRequest) error {
	resp, err := c.doPut(fmt.Sprintf("/api/users/%s/password", id), req)
	if err != nil {
		return fmt.Errorf("error ChangeUserPassword : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetCurrentUser() (*User, error) {
	resp, err := c.doGet("/api/users/current", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetCurrentUser : %w", err)
	}
	return bodyToStructure[User](resp.Body)
}

func (c *Client) GetUserInvites() ([]Invite, error) {
	resp, err := c.doGet("/api/users/invites", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetUserInvites : %w", err)
	}
	return bodyToSlice[Invite](resp.Body)
}

func (c *Client) GetInvite(nameOrEmailOrID string) (*Invite, error) {
	invites, err := c.GetUserInvites()
	if err != nil {
		return nil, err
	}
	for _, inv := range invites {
		if inv.Email == nameOrEmailOrID || inv.Name == nameOrEmailOrID || inv.ID == nameOrEmailOrID {
			return &inv, nil
		}
	}
	return nil, fmt.Errorf("invitation %s not found", nameOrEmailOrID)
}

func (c *Client) CreateUserInvite(req *CreateInviteRequest) (*Invite, error) {
	resp, err := c.doPost("/api/users/invites", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateUserInvite : %w", err)
	}
	return bodyToStructure[Invite](resp.Body)
}

func (c *Client) DeleteUserInvite(inviteID string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/users/invites/%s", inviteID))
	if err != nil {
		return fmt.Errorf("error DeleteUserInvite : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) RegenerateUserInvite(inviteID string, req *RegenerateInviteRequest) (*Invite, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/users/invites/%s/regenerate", inviteID), req)
	if err != nil {
		return nil, fmt.Errorf("error RegenerateUserInvite : %w", err)
	}
	return bodyToStructure[Invite](resp.Body)
}

func (c *Client) GetInviteInfo(token string) (*InviteInfo, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/users/invites/%s", token), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetInviteInfo : %w", err)
	}
	return bodyToStructure[InviteInfo](resp.Body)
}

func (c *Client) AcceptInvite(token string, req *AcceptInviteRequest) (*AcceptInviteResponse, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/users/invites/%s/accept", token), req)
	if err != nil {
		return nil, fmt.Errorf("error AcceptInvite : %w", err)
	}
	return bodyToStructure[AcceptInviteResponse](resp.Body)
}

func (c *Client) MergeUser(base, update *User) *User {
	if base == nil {
		return update
	}
	if update == nil {
		return base
	}

	merged := *base

	if update.ID != "" {
		merged.ID = update.ID
	}
	if update.Email != "" {
		merged.Email = update.Email
	}
	if update.Password != "" {
		merged.Password = update.Password
	}
	if update.Name != "" {
		merged.Name = update.Name
	}
	if update.Role != "" {
		merged.Role = update.Role
	}
	if update.Status != "" {
		merged.Status = update.Status
	}
	if update.LastLogin != "" {
		merged.LastLogin = update.LastLogin
	}
	if update.Issued != "" {
		merged.Issued = update.Issued
	}
	if update.IdpID != "" {
		merged.IdpID = update.IdpID
	}
	if len(update.AutoGroups) > 0 {
		merged.AutoGroups = update.AutoGroups
	}
	if update.IsCurrent != nil {
		merged.IsCurrent = update.IsCurrent
	}
	if update.IsServiceUser != nil {
		merged.IsServiceUser = update.IsServiceUser
	}
	if update.IsBlocked != nil {
		merged.IsBlocked = update.IsBlocked
	}
	if update.PendingApproval != nil {
		merged.PendingApproval = update.PendingApproval
	}
	if update.Permissions != nil {
		if merged.Permissions == nil {
			merged.Permissions = &Permissions{}
		}
		if update.Permissions.IsRestricted != nil {
			merged.Permissions.IsRestricted = update.Permissions.IsRestricted
		}
		if len(update.Permissions.Modules) > 0 {
			if merged.Permissions.Modules == nil {
				merged.Permissions.Modules = make(map[string]map[string]bool)
			}
			for module, actions := range update.Permissions.Modules {
				if merged.Permissions.Modules[module] == nil {
					merged.Permissions.Modules[module] = make(map[string]bool)
				}
				for action, val := range actions {
					merged.Permissions.Modules[module][action] = val
				}
			}
		}
	}

	return &merged
}

func (c *Client) GetUserByID(id string) (*User, error) {
	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("user with id %s not found", id)
}

func (c *Client) GetUser(nameOrIDorEmail string) (*User, error) {
	if id, ok := c.cacheLookup("user:" + nameOrIDorEmail); ok {
		return c.GetUserByID(id)
	}

	user, err := c.GetUserByEmail(nameOrIDorEmail)
	if err == nil {
		c.cacheStore("user:"+user.Name, user.ID)
		c.cacheStore("user:"+user.Email, user.ID)
		return user, nil
	}

	user, err = c.GetUserByID(nameOrIDorEmail)
	if err == nil {
		c.cacheStore("user:"+user.Name, user.ID)
		c.cacheStore("user:"+user.Email, user.ID)
		return user, nil
	}

	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Name == nameOrIDorEmail {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user %s not found", nameOrIDorEmail)
}

func (c *Client) ResolveUserID(nameOrIDorEmail string) (string, error) {
	if id, ok := c.cacheLookup("user:" + nameOrIDorEmail); ok {
		return id, nil
	}

	user, err := c.GetUser(nameOrIDorEmail)
	if err != nil {
		return "", err
	}

	_ = user // pacify unused
	c.cacheStore("user:"+user.Name, user.ID)
	c.cacheStore("user:"+user.Email, user.ID)
	c.cacheStore("user:"+user.ID, user.ID)
	return user.ID, nil
}

func (c *Client) httpResponse(code int) error {
	if code < 300 {
		return nil
	}
	if code == http.StatusNotFound {
		return fmt.Errorf("resource not found (404)")
	}
	if code == http.StatusUnauthorized {
		return fmt.Errorf("authentication required (401)")
	}
	if code == http.StatusForbidden {
		return fmt.Errorf("access denied (403)")
	}
	if code >= 500 {
		return fmt.Errorf("server error (%d)", code)
	}
	return fmt.Errorf("unexpected error (code %d)", code)
}
