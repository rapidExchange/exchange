package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"rapidEx/config"
	"rapidEx/internal/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type app struct {
	c 				config.Config
	ctx				context.Context
	redisClient		*redis.Client
}

func New() (*app, error) {
	var a app
	pwd, _ := os.Getwd()
	c, err := config.LoadConfig(pwd)
	if err != nil {
		return nil, err
	}

	a.c = c

	if err := a.setRedisConn(); err != nil {
		return nil, err
	}
	return &a, nil

}

func (a *app) setRedisConn() error {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@localhost:6379/0", a.c.RedisUser, a.c.RedisPassword))

	if err != nil {
		return err
	}

	a.redisClient = redis.NewClient(opt)

	return nil
}

func (a *app) ListenAndServe() {
	fiberApp := fiber.New()

	controllers.RegisterRoutes(fiberApp)

	log.Fatal(fiberApp.Listen(":8080"))
}