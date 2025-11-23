package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config holds the configuration for the distributed lock
type Config struct {
	Redis     RedisConfig     `mapstructure:"redis"`
	Etcd      EtcdConfig      `mapstructure:"etcd"`
	MySQL     MySQLConfig     `mapstructure:"mysql"`
	ZooKeeper ZooKeeperConfig `mapstructure:"zookeeper"`
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	Addrs    []string `mapstructure:"addrs"`
	Password string   `mapstructure:"password"`
	DB       int      `mapstructure:"db"`
}

// EtcdConfig holds etcd-specific configuration
type EtcdConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	Endpoints   []string      `mapstructure:"endpoints"`
	Username    string        `mapstructure:"username"`
	Password    string        `mapstructure:"password"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
}

// MySQLConfig holds MySQL-specific configuration
type MySQLConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
}

// ZooKeeperConfig holds ZooKeeper-specific configuration
type ZooKeeperConfig struct {
	Enabled        bool          `mapstructure:"enabled"`
	Servers        []string      `mapstructure:"servers"`
	SessionTimeout time.Duration `mapstructure:"session_timeout"`
	Prefix         string        `mapstructure:"prefix"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Set default values
	viper.SetDefault("redis.enabled", false)
	viper.SetDefault("redis.addrs", []string{"localhost:6379"})
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("etcd.enabled", false)
	viper.SetDefault("etcd.endpoints", []string{"localhost:2379"})
	viper.SetDefault("etcd.dial_timeout", "5s")

	viper.SetDefault("mysql.enabled", false)
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)

	viper.SetDefault("zookeeper.enabled", false)
	viper.SetDefault("zookeeper.servers", []string{"localhost:2181"})
	viper.SetDefault("zookeeper.session_timeout", "10s")
	viper.SetDefault("zookeeper.prefix", "/locks")

	// Read from environment variables
	viper.SetEnvPrefix("DLOCK")
	viper.AutomaticEnv()

	// Read from config file if provided
	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetConfigPath finds the config file in common locations
func GetConfigPath() string {
	// Check current directory
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}

	// Check config directory
	if _, err := os.Stat(filepath.Join("config", "config.yaml")); err == nil {
		return filepath.Join("config", "config.yaml")
	}

	// Check home directory
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "distributed-lock", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Return empty if no config file found
	return ""
}
