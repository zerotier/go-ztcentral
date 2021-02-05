package ztcentral

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-retryablehttp"
)

type Member struct {
	ID                  string       `json:"id,omitempty"`
	Type                string       `json:"type,omitempty"`
	Clock               int64        `json:"clock,omitempty"`
	NetworkID           string       `json:"networkId,omitempty"`
	NodeID              string       `json:"nodeId,omitempty"`
	ControllerID        string       `json:"controllerId,omitempty"`
	Hidden              bool         `json:"hidden"`
	Name                string       `json:"name,omitempty"`
	Online              bool         `json:"online"`
	Description         string       `json:"description,omitempty"`
	Config              MemberConfig `json:"config"`
	LastOnline          int64        `json:"lastOnline,omitempty"`
	PhysicalAddress     *string      `json:"physicalAddress,omitempty"`
	PhysicalLocation    *string      `json:"physicalLocation,omitempty"`
	ClientVersion       string       `json:"clientVersion,omitempty"`
	ProtocolVersion     int          `json:"protocolVersion"`
	SupportsRulesEngine bool         `json:"supportsRulesEngine,omitempty"`
}

type MemberConfig struct {
	ActiveBridge         bool      `json:"activeBridge"`
	Address              string    `json:"address,omitempty"`
	Authorized           bool      `json:"authorized"`
	Capabilities         []uint    `json:"capabilities,omitempty"`
	CreationTime         int64     `json:"creationTime"`
	ID                   string    `json:"id,omitempty"`
	Identity             string    `json:"identity,omitempty"`
	IPAssignments        []string  `json:"ipAssignments,omitempty"`
	LastAuthorizedTime   int64     `json:"lastAuthorizedTime"`
	LastDeauthorizedTime int64     `json:"lastDeauthorizedTime"`
	NoAutoAssignIPs      bool      `json:"noAutoAssignIps"`
	NetworkID            string    `json:"nwid,omitempty"`
	ObjectType           string    `json:"objtype,omitempty"`
	RemoteTraceLevel     int       `json:"remoteTraceLevel"`
	RemoteTraceTarget    *string   `json:"remoteTraceTarget,omitempty"`
	Revision             uint64    `json:"revision"`
	Tags                 [][2]uint `json:"tags,omitempty"`
	VersionMajor         int       `json:"vMajor"`
	VersionMinor         int       `json:"vMinor"`
	VersionRev           int       `json:"vRev"`
	VersionProtocol      int       `json:"vProto"`
}

type MemberList []Member

func (c *Client) GetMembers(ctx context.Context, networkID string) (MemberList, error) {
	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/network/%s/member", c.BaseURL, networkID), nil)
	if err != nil {
		return nil, err
	}

	res := make(MemberList, 0)
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetMember(ctx context.Context, networkID, memberID string) (*Member, error) {
	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), nil)
	if err != nil {
		return nil, err
	}

	res := Member{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateMember(ctx context.Context, m *Member) (*Member, error) {
	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, m.NetworkID, m.NodeID), reqBody)
	if err != nil {
		return nil, err
	}

	res := Member{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) CreateAuthorizedMember(ctx context.Context, networkID, memberID, name string) (*Member, error) {
	m := Member{
		ID:        fmt.Sprintf("%s-%s", networkID, memberID),
		NetworkID: networkID,
		NodeID:    memberID,
		Name:      name,
		Config: MemberConfig{
			Authorized: true,
		},
	}

	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), reqBody)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(ctx, req, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) AuthorizeMember(ctx context.Context, networkID, memberID string) (*Member, error) {
	m := Member{
		Config: MemberConfig{
			Authorized: true,
		},
	}

	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), reqBody)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(ctx, req, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) DeauthorizeMember(ctx context.Context, networkID, memberID string) (*Member, error) {
	m := Member{
		Config: MemberConfig{
			Authorized: false,
		},
	}

	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), reqBody)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(ctx, req, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) DeleteMember(ctx context.Context, m *Member) error {
	req, err := retryablehttp.NewRequest("DELETE", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, m.NetworkID, m.NodeID), nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteMemberByID(ctx context.Context, networkID, memberID string) error {
	m := Member{
		NetworkID: networkID,
		NodeID:    memberID,
	}

	return c.DeleteMember(ctx, &m)
}
