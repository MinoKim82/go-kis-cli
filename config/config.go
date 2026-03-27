package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	CfgFile string
	EnvName string // "mock" or "prod"
)

// AppConfig holds the configuration loaded from YAML or Env.
// This structure is meant to hold the profile data.
type AppConfig struct {
	AppKey    string `mapstructure:"appkey"`
	AppSecret string `mapstructure:"appsecret"`
	Cano      string `mapstructure:"cano"`
	PrdtCd    string `mapstructure:"prdt_cd"`
}

// InitConfig initializes viper to read the config file.
func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error reading home directory:", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kis-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kis-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// GetBaseURL returns the API endpoint based on the selected environment.
func GetBaseURL() string {
	if EnvName == "prod" {
		return "https://openapi.koreainvestment.com:9443"
	}
	return "https://openapivts.koreainvestment.com:29443" // Mock is default
}

// LoadProfile loads the specific environment's credentials from Viper.
// Users can define multiple profiles in ~/.kis-cli.yaml like:
// prod:
//
//	appkey: "..."
//	appsecret: "..."
//
// mock:
//
//	appkey: "..."
func LoadProfile() (*AppConfig, error) {
	var cfg AppConfig
	// Read the sub-tree corresponding to the EnvName
	sub := viper.Sub(EnvName)
	if sub == nil {
		return nil, fmt.Errorf("configuration profile '%s' not found in config file", EnvName)
	}

	err := sub.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
