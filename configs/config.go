package configs

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ServerPort string        `mapstructure:"SERVER_PORT"`
	RedisAddr  string        `mapstructure:"REDIS_ADDR"`
	IPLimit    int           `mapstructure:"IP_LIMIT"`
	TokenLimit int           `mapstructure:"TOKEN_LIMIT"`
	BlockTime  time.Duration `mapstructure:"BLOCK_TIME"`
}

func LoadConfig() (*Config, error) {
	targetFileName := ".env"
	var cfg *Config

	exeDir, _ := os.Executable()
	rootDir := filepath.Dir(exeDir)

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(rootDir)
	viper.SetConfigFile(targetFileName)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
