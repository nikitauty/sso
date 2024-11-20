package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"15m"`
	RefreshTTL     time.Duration `yaml:"refresh_ttl" env-default:"1h"`
	GRPC           GRPCConfig    `yaml:"grpc"`
	PostgresConfig `yaml:"postgres"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true" env-default:"5432"`
	Username string `yaml:"username" env-required:"true" env-default:"postgres"`
	Password string `yaml:"password" env-required:"true" env-default:"postgres"`
	Database string `yaml:"database" env-required:"true"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config file path is required")
	}

	return LoadByPath(configPath)
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		os.Getenv("CONFIG_PATH")
	}
	return res
}

func LoadByPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
