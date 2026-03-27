package client_test

import (
	"testing"

	"github.com/MinoKim82/go-kis-cli/config"
	"github.com/MinoKim82/go-kis-cli/pkg/client"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupMockConfig() {
	config.EnvName = "test"
	viper.Set("test.appkey", "test-key")
	viper.Set("test.appsecret", "test-secret")
	viper.Set("test.cano", "12345678")
	viper.Set("test.prdt_cd", "01")
}

func TestKISClient_Request(t *testing.T) {
	setupMockConfig()

	c, err := client.NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)

	req := c.Request("TEST_TR_ID", "mock-token")
	
	// Verify headers are correctly injected
	assert.Equal(t, "Bearer mock-token", req.Header.Get("authorization"))
	assert.Equal(t, "test-key", req.Header.Get("appkey"))
	assert.Equal(t, "test-secret", req.Header.Get("appsecret"))
	assert.Equal(t, "TEST_TR_ID", req.Header.Get("tr_id"))
	assert.Equal(t, "P", req.Header.Get("custtype"))
}

func TestIsSuccess(t *testing.T) {
	assert.True(t, client.IsSuccess("0"))
	assert.False(t, client.IsSuccess("-1"))
	assert.False(t, client.IsSuccess("1"))
}
