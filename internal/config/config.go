package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8080"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	EnvLocal                      = "local"
	Prod                          = "prod"
)

type (
	Config struct {
		Environment string
		Postgres    PostgresConfig
		HTTP        HTTPConfig
	}
	PostgresConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}
	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func setFromEnv(cfg *Config) {
	// TODO use envconfig https://github.com/kelseyhightower/envconfig
	cfg.Postgres.URI = os.Getenv("POSTGRES_URI")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASS")
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.Environment = os.Getenv("APP_ENV")
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()

}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	return nil
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
}
