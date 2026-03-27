package client

import (
	"fmt"
	"net/http"

	"github.com/MinoKim82/go-kis-cli/config"
	"github.com/go-resty/resty/v2"
)

// KISClient wraps the go-resty client
type KISClient struct {
	RestyClient *resty.Client
	Profile     *config.AppConfig
}

// NewClient creates a new KIS API client
func NewClient() (*KISClient, error) {
	profile, err := config.LoadProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to load profile: %w", err)
	}

	rc := resty.New()
	// Allow httpmock to intercept resty
	if config.EnvName == "test" || config.EnvName == "mock" {
		rc.SetTransport(http.DefaultTransport)
	}
	rc.SetBaseURL(config.GetBaseURL())
	rc.SetHeader("Content-Type", "application/json; charset=utf-8")

	return &KISClient{
		RestyClient: rc,
		Profile:     profile,
	}, nil
}

// Request is a wrapper around the resty request, automatically injecting required headers.
// trID is the transaction ID representing the KIS API endpoint behavior.
func (c *KISClient) Request(trID string, token string) *resty.Request {
	req := c.RestyClient.R().
		SetHeader("authorization", fmt.Sprintf("Bearer %s", token)).
		SetHeader("appkey", c.Profile.AppKey).
		SetHeader("appsecret", c.Profile.AppSecret).
		SetHeader("tr_id", trID)

	req.SetHeader("custtype", "P")
	return req
}

// KISError represents a standard error response from the Korea Investment REST API.
type KISError struct {
	RtCd  string `json:"rt_cd"`
	MsgCd string `json:"msg_cd"`
	Msg1  string `json:"msg1"`
}

func (e *KISError) Error() string {
	return fmt.Sprintf("KIS API Error [%s]: %s", e.MsgCd, e.Msg1)
}

// IsSuccess checks the standard rt_cd field from a generic map response.
// KIS API returns rt_cd == "0" on success.
func IsSuccess(rtCd string) bool {
	return rtCd == "0"
}
