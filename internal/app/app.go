package app

import (
	"context"
	"fmt"
	"rapidEx/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type app struct {
	c 				config.Config
	ctx				context.Context
	redisClient		*redis.Client
}

func New() (*app, error) {
	var c config.Config

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	var a app

	if err := a.setRedisConn(); err != nil {
		return nil, err
	}

	return &a, nil

}

func (a *app) setRedisConn() error {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@localhost:6379", a.c.RedisUser, a.c.RedisPassword))

	if err != nil {
		return err
	}

	a.redisClient = redis.NewClient(opt)

	return nil
} 