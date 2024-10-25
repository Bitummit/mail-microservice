package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)


type Config struct {
	Env string	`yaml:"env" env-default:"dev"`
	GrpcServer	`yaml:"grpc_server"`
	Kafka 		`yaml:"kafka_broker"`
	Email		`yaml:"email"`

}

type GrpcServer struct {
	GrpcAddress string `yaml:"grpc_address" env-default:"0.0.0.0:5400"`
}

type Kafka struct {
	KafkaAddress string `yaml:"kafka_address" env-default:"0.0.0.0:9092"`
}

type Email struct {
	Server string 	`yaml:"smtp_server" env-default:"smtp.mail.ru"`
	Port int		`yaml:"smtp_port" env-default:"587"`
}


func InitConfig() *Config{
	
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file!")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("Empty path")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalln("Can not find config file")
	} 
	
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalln("Error in reading config file!", err)
	}
	
	return &cfg
	
}