package config

import (
	"fmt"
	"os"

	
	"github.com/spf13/viper"
	"github.com/yusufaniki/muslim_tech/pkg/logger"
)





type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Pass     string `mapstructure:"DB_PASS"`
	Name     string `mapstructure:"DB_NAME"`
	Schema   string `mapstructure:"DB_SCHEMA"`
}


type RedisConfig struct {
	Host string    `mapstructure:"REDIS_HOST"`
	Port string    `mapstructure:"REDIS_PORT"`
	Pass string    `mapstructure:"REDIS_PASS"`
	User string    `mapstructure:"REDIS_USER"`
	DB  int        `mapstructure:"REDIS_DB"`
	Enabled bool   `mapstructure:"REDIS_ENABLED"`
}

type MailConfig struct {
	Host string `mapstructure:"MAIL_HOST"`
	Port string `mapstructure:"MAIL_PORT"`
	User string `mapstructure:"MAIL_USER"`
	Pass string `mapstructure:"MAIL_PASS"`
	From string `mapstructure:"MAIL_FROM"`
}
type JWTConfig struct {
	Secret    string   `mapstructure:"JWT_SECRET"`
	Duration  int      `mapstructure:"JWT_DURATION"`
}
type Config struct {  
	DBSource       string          `mapstructure:"DB_SOURCE"`
	DBSourceLocal  string          `mapstructure:"DB_SOURCE_LOCAL"`
	Port           string          `mapstructure:"APP_PORT"`
	AppEnv         string          `mapstructure:"APP_ENV"`
	EnvPath        string          `mapstructure:"ENV_PATH"`
	JWT            JWTConfig       `mapstructure:",squash"`
	Database       DatabaseConfig  `mapstructure:",squash"`
	Redis          RedisConfig     `mapstructure:",squash"`
    Mail		   MailConfig      `mapstructure:",squash"`
    JWTSecret	  string          `mapstructure:"JWT_SECRET"`
}


func LoadConfig() (Config, error) {
	configPath := os.Getenv("ENV_PATH")
	if	configPath == "" {
		configPath = "./"
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigFile(configPath + ".env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
     
	appLog := logger.CreateZapLogger()


	var config Config
    if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			appLog.Warn("No configuration file found, using default values and environment variables")
			return config, err
		} else {
			appLog.Error("Error reading configuration file.", map[string]interface{}{"error": err})
			return config, fmt.Errorf("error reading configuration file: %v", err)
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		appLog.Error("Unable to decode into struct.", map[string]interface{}{"error": err})
		return config, fmt.Errorf("unable to decode into struct: %v", err)
	}

	if err := ValidateConfig(config); err != nil {
		return config, err
	}
	
	return config, nil 
}


func ValidateConfig(config Config) error {
	if config.DBSource == "" {
		return fmt.Errorf("DB_SOURCE is required")
	}
	
	if config.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	return nil
}
