package redisconnect

import (
	"fmt"
	"os"
	"rapidEx/config"

	"github.com/redis/go-redis/v9"
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
