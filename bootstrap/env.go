package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ApiPort                string `mapstructure:"API_PORT"`
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBDebug                string `mapstructure:"DB_DEBUG"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	DBName                 string `mapstructure:"DB_NAME"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPort                 string `mapstructure:"DB_PORT"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
