package auth_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/MinoKim82/go-kis-cli/config"
	"github.com/MinoKim82/go-kis-cli/pkg/auth"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupConfig() {
	config.EnvName = "test"
	viper.Set("test.appkey", "mock-appkey")
	viper.Set("test.appsecret", "mock-appsecret")
}

func TestIssueToken(t *testing.T) {
	setupConfig()

	// httpmock is active globally


	mockResp := auth.TokenResponse{
		AccessToken: "mock-access-token-123",
		TokenType:   "Bearer",
		ExpiresIn:   86400,
	}

	httpmock.RegisterResponder("POST", "https://openapivts.koreainvestment.com:29443/oauth2/tokenP",
		func(req *http.Request) (*http.Response, error) {
			
			var body map[string]string
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, "Bad Request"), nil
			}
			
			assert.Equal(t, "client_credentials", body["grant_type"])
			assert.Equal(t, "mock-appkey", body["appkey"])
			assert.Equal(t, "mock-appsecret", body["appsecret"])
			
			resp, err := httpmock.NewJsonResponse(200, mockResp)
			return resp, err
		},
	)

	// In test, running IssueToken might be flaky if it uses its own client. 
	// Wait, IssueToken inside auth.go creates a NEW client. So mocking c.RestyClient.GetClient() here won't work 
	// because IssueToken calls client.NewClient() internally.
	// That means we need `httpmock.Activate()` to globals if resty uses default transport.
	// But resty uses its own transport by default unless we hook globally.
	// Let's activate globally.
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// issueToken should hit the mock
	err := auth.IssueToken()
	assert.NoError(t, err)

	// Since we cache locally to ~/.kis-cli-token.json, we can check GetValidToken
	token, err := auth.GetValidToken()
	assert.NoError(t, err)
	assert.Equal(t, "mock-access-token-123", token)
}
