package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Server   ServerConf
	Storage  StorageConf
}

type LoggerConf struct {
	Level string
	File  string
}

type DatabaseConf struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSL      string
}

type ServerConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type string
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	}

	return &config, nil
}

func (d *DatabaseConf) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSL)
}

func (s *ServerConf) DSN() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
