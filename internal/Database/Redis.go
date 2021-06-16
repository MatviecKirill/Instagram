package redisDB

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/url"
	"os"
	"strings"
)

var ctx = context.Background()

func Init() error {
	herokuURL := os.Getenv("REDIS_URL")
	password := ""
	if !strings.Contains(herokuURL, "localhost") && herokuURL != "" {
		if parsedURL, err := url.Parse(herokuURL); err != nil {
			return errors.New("can't parse heroku redis URL")
		} else {
			password, _ = parsedURL.User.Password()
			herokuURL = parsedURL.Host
		}
	} else {
		return errors.New("heroku redis URL not found")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     herokuURL,
		Password: password,
		DB:       0,
	})

	fmt.Println("Redis connect successfully")
	rdb.FlushAll(ctx)
	/*err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}*/
	return nil
}
