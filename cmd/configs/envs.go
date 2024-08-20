package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	Port           int    `mapstructure:"PORT"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`

	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPwd  string `mapstructure:"DB_PWD"`
	DBName string `mapstructure:"DB_NAME"`

	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

var Envs = initConfig()

func initConfig() Env {

	env := Env{}
	viper.SetConfigFile("development.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the environment file : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return env
}
