package redis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (rdb *redis.Client) {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	addr := redisHost + ":" + redisPort

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       1,
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis: ", err.Error())
	} else {
		log.Println("Connected to Redis: " + addr + " - " + status)
	}

	return rdb
}
