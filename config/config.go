package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//App config
type App struct {
	//Server config
	Server Server
	//Db has the database config
	Db Database
}

//Server is the server config
type Server struct {
	//Port on which the app is running
	Port string
}

var appConfig App

func GetConfig() App {
	return appConfig
}

func init() {
	/*
	 * We will specify the config
	 * Then will load from the config
	 * Then parse the config
	 */
	log.Info("going to read the config from config.json file and environment variables")
	//specify the config file and it paths
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	//loading the config
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error while reading the config file")
	}

	//parsing the config
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error while parsing config")
	}

	log.Info("loaded the config file")
}
