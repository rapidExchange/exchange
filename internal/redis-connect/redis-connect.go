package redisconnect

import (
	"os"
	"fmt"
	"github.com/redis/go-redis/v9"
	"rapidEx/config"
)

func SetRedisConn() (*redis.Client, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	c, err := config.LoadConfig(pwd)

	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@localhost:6379/1", c.RedisUser, c.RedisPassword))
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}
