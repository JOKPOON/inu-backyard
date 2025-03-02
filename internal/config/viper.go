package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func SetConfig(class interface{}) {
	//Get the config path from the environment
	name := os.Getenv("CONFIG_NAME")
	if name == "" {
		name = "config"
	}
	viper.SetConfigName(name)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&class)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

// TODO: read from struct not viper
func PrintConfig() {
	for _, field := range viper.AllKeys() {
		if viper.IsSet(field) {
			log.Printf("%s: %v", field, viper.Get(field))
		}
	}
}
