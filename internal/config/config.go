// Package config provides loading and parsing of YAML configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	OpenCode struct {
		Host    string `mapstructure:"host"`
		Port    int    `mapstructure:"port"`
		Timeout int    `mapstructure:"timeout"`
	} `mapstructure:"opencode"`

	Generation struct {
		Style string `mapstructure:"style"`
		Model struct {
			Provider string `mapstructure:"provider"`
			ModelID  string `mapstructure:"model_id"`
		} `mapstructure:"model"`
	} `mapstructure:"generation"`

	Cache struct {
		Enabled  bool   `mapstructure:"enabled"`
		TTL      string `mapstructure:"ttl"`
		Location string `mapstructure:"location"`
	} `mapstructure:"cache"`

	Git struct {
		StagedOnly bool `mapstructure:"staged_only"`
	} `mapstructure:"git"`
}

var cfg *Config

func Initialize(cfgFile string) error {
	viper.SetDefault("opencode.host", "localhost")
	viper.SetDefault("opencode.port", 4096)
	viper.SetDefault("opencode.timeout", 120)

	viper.SetDefault("generation.style", "conventional")
	viper.SetDefault("generation.model.provider", "google")
	viper.SetDefault("generation.model.model_id", "antigravity-gemini-3-flash")

	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", "24h")

	viper.SetDefault("git.staged_only", true)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(homeDir, ".config", "commit-gen"))
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			err := SaveConfig()
			if err != nil {
				fmt.Printf("Warning: failed to save config: %v\n", err)
			}
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	viper.SetEnvPrefix("COMMIT_GEN")
	viper.AutomaticEnv()

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func Get() *Config {
	if cfg == nil {
		err := Initialize("")
		if err != nil {
			fmt.Printf("Warning: failed to initialize config: %v\n", err)
		}
	}
	return cfg
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func SaveConfig() error {
	return viper.WriteConfig()
}
