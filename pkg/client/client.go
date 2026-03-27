package client

import (
	"fmt"

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
	rc.SetBaseURL(config.GetBaseURL())
	rc.SetHeader("Content-Type", "application/json; charset=utf-8")

	return &KISClient{
		RestyClient: rc,
		Profile:     profile,
	}, nil
}
