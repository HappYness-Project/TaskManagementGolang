package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	Port           int    `mapstructure:"PORT"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
	LogLevel       string `mapstructure:"LOG_LEVEL`

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
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		fmt.Println(e.Name())
	}

	workingdir, _ := os.Getwd()
	fmt.Println("Current Dir: " + workingdir)
	if envString == "" {
		viper.SetConfigFile(workingdir + "/../dev-env/dev.env")
	} else if envString == "development" {
		entries, err := os.ReadDir(workingdir + "/../")
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range entries {
			fmt.Println(e.Name())
		}

		viper.SetConfigFile("../.env")
	}
	env := Env{}
	err = viper.ReadInConfig()
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
