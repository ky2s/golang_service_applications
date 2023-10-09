package main

import (
	"fmt"
	"snapin-form/config"
	"snapin-form/routes"
	"strconv"

	"github.com/spf13/viper"
)

func main() {

	// Set the file name of the configurations file
	viper.SetConfigName("gifnoc")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	var conf config.Configurations
	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	r := routes.SetupRoutes(conf)

	r.Run(":" + strconv.Itoa(conf.Server.Port))
	// r.RunTLS(conf.Server.Hostname+":"+strconv.Itoa(conf.Server.Ssl_Port), conf.SSL_PUBLIC_KEY, conf.SSL_PRIVATE_KEY)
}
