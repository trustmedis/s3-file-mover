package lib

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	ACCESS_KEY    string   `validate:"required"`
	ACCESS_SECRET string   `validate:"required"`
	REGION        string   `validate:"required"`
	BUCKET        string   `validate:"required"`
	ENDPOINT      string   `validate:"required"`
	WATCH_DIR     []string `validate:"required"`
	AUTO_CLEANUP  bool
}

func LoadConfig() *Config {
	configLocation := flag.String("config", "./config.toml", "config.toml location")
	flag.Parse()
	log.Println("Loading config from", *configLocation)

	var config Config
	_, err := toml.DecodeFile(*configLocation, &config)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Println(err.Tag(), err.Field())
		}
		os.Exit(1)
	}

	return &config
}
