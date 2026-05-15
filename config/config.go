package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Postgres PostgresConfig `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
	JWT      JWTConfig      `mapstructure:",squash"`
}

type AppConfig struct {
	Env    string `mapstructure:"APP_ENV"`
	Port   string `mapstructure:"APP_PORT"`
	Secret string `mapstructure:"APP_SECRET_KEY"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type JWTConfig struct {
	AccessSecret        string `mapstructure:"JWT_ACCESS_SECRET"`
	RefreshSecret       string `mapstructure:"JWT_REFRESH_SECRET"`
	AccessExpireMinutes int    `mapstructure:"JWT_ACCESS_EXPIRE_MINUTES"`
	RefreshExpireDays   int    `mapstructure:"JWT_REFRESH_EXPIRE_DAYS"`
}

func Load() (*Config, error) {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("env file not found, using system env")
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.Name,
		p.SSLMode,
	)
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
