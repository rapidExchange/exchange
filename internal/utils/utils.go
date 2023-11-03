package utils
import (
	"strconv"
	"os"
	"fmt"
	"github.com/redis/go-redis/v9"
	"rapidEx/config"
)

func MapFloatToString(m map[float64]float64) map[string]string {
	sMap := make(map[string]string)

	for k, v := range m {
		sKey := strconv.FormatFloat(k, 'f', -1, 64)
		sVal := strconv.FormatFloat(v, 'f', -1, 64)
		sMap[sKey] = sVal
	}
	return sMap
}

func MapStringToFloat(m map[string]string) (map[float64]float64, error) {
	fMap := make(map[float64]float64)

	for k, v := range m {
		fKey, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil, err
		}
		fVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}

		fMap[fKey] = fVal
	}

	return fMap, nil
}

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
