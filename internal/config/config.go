package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	
	// Support environment-specific configs
	// Priority: config.local.yaml > config.yaml
	// Use config.local.yaml for local development/debugging
	env := os.Getenv("APP_ENV")
	if env == "" {
		// Check if local config exists for development
		if _, err := os.Stat(path + "/config.local.yaml"); err == nil {
			viper.SetConfigName("config.local")
		} else {
			viper.SetConfigName("config")
		}
	} else {
		viper.SetConfigName("config." + env)
	}
	
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".","_"))

	err = viper.ReadInConfig()
	if err!=nil{
		if _,ok:=err.(viper.ConfigFileNotFoundError);ok{
			return
		}
	}
	 
	err = viper.Unmarshal(&config)
	return
}
