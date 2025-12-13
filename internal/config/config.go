package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	JWTSecret     string        `mapstructure:"JWT_SECRET"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`

	// Redis
    RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
    RedisPassword string `mapstructure:"REDIS_PASSWORD"`
    RedisDb       int    `mapstructure:"REDIS_DB"`

    // Email (SMTP)
    EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
    EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
    EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`

	CloudinaryURL string `mapstructure:"CLOUDINARY_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  
	viper.SetConfigName("app") 
	viper.SetConfigType("env") 

	viper.AutomaticEnv() 
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	
	err = viper.Unmarshal(&config)
	return
}