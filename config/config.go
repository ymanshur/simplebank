package config

import (
	"time"

	"github.com/spf13/viper"
	"github.com/ymanshur/simplebank/pkg/util"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Debug                    bool          `mapstructure:"DEBUG"`
	Environment              string        `mapstructure:"ENVIRONMENT"`
	AllowedOrigins           []string      `mapstructure:"ALLOWED_ORIGINS"`
	DBSource                 string        `mapstructure:"DB_SOURCE"`
	DBMigrationURL           string        `mapstructure:"DB_MIGRATION_URL"`
	RedisAddress             string        `mapstructure:"REDIS_ADDRESS"`
	RedisUsername            string        `mapstructure:"REDIS_USERNAME"`
	RedisPassword            string        `mapstructure:"REDIS_PASSWORD"`
	HTTPServerAddress        string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress        string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	GRPCGatewayServerAddress string        `mapstructure:"GRPC_GATEWAY_SERVER_ADDRESS"`
	TokenSymmetricKey        string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration      time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration     time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName          string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress       string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword      string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

// LoadConfig reads configuration from file or environment variables.
// It will read file in root of project.
func LoadConfig() (Config, error) {
	// projectRoot is default config path
	projectRoot := util.RootDir()

	viper.AddConfigPath(projectRoot)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, yml, etc.

	// AutomaticEnv will override config file
	viper.AutomaticEnv()

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
