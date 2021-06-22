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
var client *redis.Client

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

	client = redis.NewClient(&redis.Options{
		Addr:     herokuURL,
		Password: password,
		DB:       0,
	})
	fmt.Println("Redis connect successfully")
	return nil
}

func Exist(key string) bool {
	if exist, err := client.Exists(ctx, key).Result(); err == nil && exist == 1 {
		return true
	}
	return false
}

func Set(key, value string) {
	client.Del(ctx, key)
	client.Set(ctx, key, value, 0)
}

func Get(key string) string {
	if value, err := client.Get(ctx, key).Result(); err == nil {
		return value
	} else {
		fmt.Println(err)
		return ""
	}
}

func Del(key string) {
	client.Del(ctx, key)
}

func SAdd(key string, values interface{}) {
	client.Del(ctx, key)
	client.SAdd(ctx, key, values)
}

func SMembers(key string) []string {
	if values, err := client.SMembers(ctx, key).Result(); err == nil {
		return values
	} else {
		fmt.Println(err)
		return nil
	}
}
