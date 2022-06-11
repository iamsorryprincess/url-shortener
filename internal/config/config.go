package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Address            string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL            string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath        string `env:"FILE_STORAGE_PATH"`
	DBConnectionString string `env:"DATABASE_DSN" envDefault:""`
	WorkersCount       int    `env:"WORKERS_COUNT" envDefault:"1"`
	WorkerPoolSize     int    `env:"WORKER_POOL_SIZE" envDefault:"1"`
}

func ParseConfiguration() (*Configuration, error) {
	configuration := &Configuration{}
	err := env.Parse(configuration)

	if err != nil {
		return nil, err
	}

	flag.StringVar(&configuration.Address, "a", configuration.Address, "server address")
	flag.StringVar(&configuration.BaseURL, "b", configuration.BaseURL, "base url")
	flag.StringVar(&configuration.StoragePath, "f", configuration.StoragePath, "file storage path")
	flag.StringVar(&configuration.DBConnectionString, "d", configuration.DBConnectionString, "db connection string")
	flag.Parse()
	return configuration, nil
}
