package service

import (
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
)

/* vars */
var redisClient *redis.Client

/* consts */
const defaultCartCapacity = 5

/* Redis consts */
const defaultRedisAddress = "redis:6379"
const defaultRedisPassword = ""
const defaultRedisDB = 0

var redisAddress = func() string {
	val, found := os.LookupEnv("REDIS_ADDRESS")
	if !found {
		val = defaultRedisAddress
	}
	return val
}()

var redisPassword = func() string {
	val, found := os.LookupEnv("REDIS_PASSWORD")
	if !found {
		val = defaultRedisPassword
	}
	return val
}()

var redisDb = func() int {
	valStr, found := os.LookupEnv("REDIS_DB")
	var val int
	if !found {
		val = defaultRedisDB
	} else {
		var err error
		val, err = strconv.Atoi(valStr)
		if err != nil {
			log.Println("Unable to parse REDIS_DB: ", err)
			log.Println("Falling back to default value ", defaultRedisDB)
			val = defaultRedisDB
		}
	}
	return val
}()
