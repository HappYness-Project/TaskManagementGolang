package configs

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	Port           string `mapstructure:"PORT"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
	LogLevel       string `mapstructure:"LOG_LEVEL`
	Host           string `mapstructure:"HOST"`

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

var AccessToken string // updated from the main package.

func InitConfig(envString string) Env {
	env := Env{}
	if envString == "" {
		workingdir, _ := os.Getwd()
		viper.SetConfigFile(workingdir + "/../dev-env/local.env")
	} else if envString == "local" {
		viper.SetConfigFile("local.env")
	} else if envString == "development" {
		env.AppEnv = envString
		env.Host = "0.0.0.0"
		env.Port = os.Getenv("PORT")
		env.DBHost = os.Getenv("DB_HOST")
		env.DBName = os.Getenv("DB_NAME")
		env.DBPort = os.Getenv("DB_PORT")
		env.DBUser = os.Getenv("DB_USER")
		env.DBPwd = os.Getenv("DB_PWD")
		env.AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
		env.RefreshTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
		return env
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the environment file : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}
	if envString == "" {
		env.DBHost = "localhost"
		env.DBPort = "8010"
	}
	return env
}
