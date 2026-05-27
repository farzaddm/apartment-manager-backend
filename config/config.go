package config

import (
	"fmt"

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

	viper.SetConfigName("prev.env") // WARNING !!!!
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	viper.BindEnv("APP_ENV")
	viper.BindEnv("APP_PORT")
	viper.BindEnv("APP_SECRET_KEY")

	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")

	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("REDIS_DB")

	viper.BindEnv("JWT_ACCESS_SECRET")
	viper.BindEnv("JWT_REFRESH_SECRET")
	viper.BindEnv("JWT_ACCESS_EXPIRE_MINUTES")
	viper.BindEnv("JWT_REFRESH_EXPIRE_DAYS")

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (p PostgresConfig) DSN() string {
	sslMode := p.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Name, sslMode,
	)
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
