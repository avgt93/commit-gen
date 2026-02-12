package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

/**
 * Config holds all configuration settings for commit-gen.
 */
type Config struct {
	OpenCode struct {
		Mode    string `mapstructure:"mode"`
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
		StagedOnly  bool   `mapstructure:"staged_only"`
		Editor      string `mapstructure:"editor"`
		MaxDiffSize int    `mapstructure:"max_diff_size"`
	} `mapstructure:"git"`
}

var cfg *Config

/**
 * Initialize loads and parses the configuration from file, environment, and defaults.
 *
 * @param cfgFile - Path to a specific config file, or empty for default locations
 * @returns An error if config loading fails
 */
func Initialize(cfgFile string) error {
	viper.SetDefault("opencode.mode", "run")
	viper.SetDefault("opencode.host", "localhost")
	viper.SetDefault("opencode.port", 4096)
	viper.SetDefault("opencode.timeout", 120)

	viper.SetDefault("generation.style", "conventional")
	viper.SetDefault("generation.model.provider", "google")
	viper.SetDefault("generation.model.model_id", "antigravity-gemini-3-pro")

	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", "24h")

	viper.SetDefault("git.staged_only", true)
	viper.SetDefault("git.editor", "cat")
	viper.SetDefault("git.max_diff_size", 32*1024)

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(homeDir, ".config", "commit-gen"))
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			_ = SaveConfig()
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

/**
 * Get returns the current configuration, initializing it if necessary.
 *
 * @returns The current Config instance
 */
func Get() *Config {
	if cfg == nil {
		err := Initialize("")
		if err != nil {
			fmt.Printf("Warning: failed to initialize config: %v\n", err)
		}
	}
	return cfg
}

/**
 * GetString retrieves a string value from the configuration.
 *
 * @param key - The configuration key to retrieve
 * @returns The string value for the given key
 */
func GetString(key string) string {
	return viper.GetString(key)
}

/**
 * GetInt retrieves an integer value from the configuration.
 *
 * @param key - The configuration key to retrieve
 * @returns The integer value for the given key
 */
func GetInt(key string) int {
	return viper.GetInt(key)
}

/**
 * GetBool retrieves a boolean value from the configuration.
 *
 * @param key - The configuration key to retrieve
 * @returns The boolean value for the given key
 */
func GetBool(key string) bool {
	return viper.GetBool(key)
}

/**
 * Set sets a configuration value.
 *
 * @param key - The configuration key to set
 * @param value - The value to set for the key
 */
func Set(key string, value interface{}) {
	viper.Set(key, value)
}

/**
 * SaveConfig writes the current configuration to file.
 *
 * @returns An error if writing fails
 */
func SaveConfig() error {
	return viper.WriteConfig()
}
