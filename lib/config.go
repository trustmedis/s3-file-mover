package lib

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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
	log.Printf("Loading config from %s. Now checking S3 credentials...", *configLocation)

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

	// Check S3 credentials
	CheckS3Credentials(&config)

	return &config
}

func CheckS3Credentials(config *Config) {
	session, err := CreateSession(config)
	if err != nil {
		log.Fatalln(err)
	}

	// Try fetching bucket information as a simple operation
	svc := s3.New(session)
	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(config.BUCKET),
	})

	if err != nil {
		log.Fatalln("Invalid S3 credentials or permissions.")
	}

	log.Println("S3 credentials and permissions are valid.")
}
