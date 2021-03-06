package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/rushuinet/version_cache"
	"log"
	"strconv"
	"time"
)

type Info struct {
	Id         int
	Key        string
	CreateTime int64
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	if pong != "PONG" || err != nil {
		log.Fatal("redis conn error")
	}
	config := &version_cache.Option{
		Redis:               client,
		Key:                 "test",
		VersionGenerateTime: 60,
		CheckTime:           5,
		KeepVersionNum:      6,
	}
	cache := version_cache.New(config)
	cache.Generate(func(key string) {
		for i := 0; i < 5; i++ {
			info := &Info{
				Id:         i,
				Key:        key,
				CreateTime: time.Now().Unix(),
			}
			data, _ := json.Marshal(info)
			cache.Redis.HSet(key, strconv.Itoa(i), data)
			time.Sleep(1 * time.Second)
		}
	})

	select {}
}
