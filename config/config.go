package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppPort string `mapstructure:"APP_PORT"`
	RedisUser     string `mapstructure:"REDIS_USER"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPassword string `mapstructure:"REDIS_PASS"`
	MysqlUser     string `mapstructure:"MYSQL_USER"`
	MysqlHost     string `mapstructure:"MYSQL_HOST"`
	MysqlPassword string `mapstructure:"MYSQL_PASSWORD"`
	MysqlDBName   string `mapstructure:"MYSQL_DBNAME"`
}

func LoadConfig(path string) (c Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
