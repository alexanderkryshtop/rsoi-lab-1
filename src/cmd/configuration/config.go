package configuration

import (
	"context"
	"flag"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger      LoggerConfig      `mapstructure:"logger"`
	Application ApplicationConfig `mapstructure:"application"`
	HTTPServer  HTTPServerConfig  `mapstructure:"server"`
}

type LoggerConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"outputPaths"`
	ErrorOutputPaths []string `mapstructure:"errorOutputPaths"`
	EnableStackTrace bool     `mapstructure:"enableStackTrace"`
}

type ApplicationConfig struct {
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

type HTTPServerConfig struct {
	Port uint64 `mapstructure:"port"`
}

func newDefaultConfig() *Config {
	return &Config{
		Logger: LoggerConfig{
			Level:            "info",
			Encoding:         "console",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EnableStackTrace: true,
		},
		Application: ApplicationConfig{
			ShutdownTimeout: 15 * time.Second,
		},
		HTTPServer: HTTPServerConfig{
			Port: 0,
		},
	}
}

func LoadConfig(ctx context.Context) (*Config, error) {
	config := newDefaultConfig()

	configPath := new(string)
	flag.StringVar(configPath, "c", "", "config path")
	flag.Parse()

	viper.SetConfigFile(*configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
