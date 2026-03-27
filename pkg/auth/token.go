package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MinoKim82/go-kis-cli/pkg/client"
)

type TokenResponse struct {
	AccessToken             string `json:"access_token"`
	AccessTokenTokenExpired string `json:"access_token_token_expired"`
	TokenType               string `json:"token_type"`
	ExpiresIn               int    `json:"expires_in"`
}

type CachedToken struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// IssueToken requests a new access token from KIS API
func IssueToken() error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}

	body := map[string]string{
		"grant_type": "client_credentials",
		"appkey":     c.Profile.AppKey,
		"appsecret":  c.Profile.AppSecret,
	}

	var result TokenResponse

	resp, err := c.RestyClient.R().
		SetBody(body).
		SetResult(&result).
		Post("/oauth2/tokenP")

	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API returned error %d: %s", resp.StatusCode(), resp.String())
	}

	fmt.Println("Token issued successfully. Caching...")

	// Cache the token
	err = cacheToken(&result)
	if err != nil {
		return fmt.Errorf("failed to cache token: %w", err)
	}
	return nil
}

func cacheToken(tr *TokenResponse) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cachePath := filepath.Join(home, ".kis-cli-token.json")

	// The token usually expires in 24 hours, let's subtract 1 hour for safety margin
	expiresAt := time.Now().Add(time.Duration(tr.ExpiresIn-3600) * time.Second)

	ct := CachedToken{
		AccessToken: tr.AccessToken,
		ExpiresAt:   expiresAt,
	}

	data, err := json.MarshalIndent(ct, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0600)
}

// GetValidToken retrieves the cached token if it's still valid, otherwise it returns empty string or issues a new one
func GetValidToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cachePath := filepath.Join(home, ".kis-cli-token.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no cached token found. Please run 'kis-cli auth login'")
		}
		return "", err
	}

	var ct CachedToken
	if err := json.Unmarshal(data, &ct); err != nil {
		return "", err
	}

	if time.Now().After(ct.ExpiresAt) {
		return "", fmt.Errorf("token expired. Please run 'kis-cli auth login' again")
	}

	return ct.AccessToken, nil
}
