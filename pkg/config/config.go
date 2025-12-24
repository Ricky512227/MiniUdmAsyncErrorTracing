package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Kubernetes KubernetesConfig `mapstructure:"kubernetes"`
	Paths      PathsConfig      `mapstructure:"paths"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Symptom    SymptomConfig    `mapstructure:"symptom"`
	Patch      PatchConfig      `mapstructure:"patch"`
}

// KubernetesConfig holds Kubernetes-related configuration
type KubernetesConfig struct {
	Namespace      string        `mapstructure:"namespace"`
	KubeconfigPath string        `mapstructure:"kubeconfig_path"`
	Timeout        time.Duration `mapstructure:"timeout"`
}

// PathsConfig holds path-related configuration
type PathsConfig struct {
	TcnVolPath string `mapstructure:"tcn_vol_path"`
	Lib64Path  string `mapstructure:"lib64_path"`
	LogPaths   []string `mapstructure:"log_paths"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// SymptomConfig holds symptom collection configuration
type SymptomConfig struct {
	ErrorKeywords    []string      `mapstructure:"error_keywords"`
	CheckInterval    time.Duration `mapstructure:"check_interval"`
	CollectionTimeout time.Duration `mapstructure:"collection_timeout"`
}

// PatchConfig holds patch application configuration
type PatchConfig struct {
	BackupEnabled bool          `mapstructure:"backup_enabled"`
	HealthTimeout time.Duration `mapstructure:"health_timeout"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	viper.SetConfigType("yaml")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Set default config paths
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("$HOME/.miniumd")
		viper.AddConfigPath("/etc/miniumd")
	}

	// Environment variables
	viper.SetEnvPrefix("MINIUDM")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Kubernetes defaults
	viper.SetDefault("kubernetes.namespace", "default")
	viper.SetDefault("kubernetes.timeout", "30s")

	// Path defaults
	viper.SetDefault("paths.tcn_vol_path", "/tcnVol")
	viper.SetDefault("paths.lib64_path", "/opt/SMAW/INTP/lib64")
	viper.SetDefault("paths.log_paths", []string{
		"/cmconfig.log",
		"/logstore/TspCore",
		"/RTPTraceError",
		"/Envoy",
		"/dumplog",
	})

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Symptom defaults
	viper.SetDefault("symptom.error_keywords", []string{
		"error", "ERROR", "fatal", "FATAL", "exception", "EXCEPTION", "panic", "PANIC",
	})
	viper.SetDefault("symptom.check_interval", "1s")
	viper.SetDefault("symptom.collection_timeout", "10m")

	// Patch defaults
	viper.SetDefault("patch.backup_enabled", true)
	viper.SetDefault("patch.health_timeout", "30s")
}

// GetHomeDir returns the home directory for configuration files
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".miniumd"), nil
}

