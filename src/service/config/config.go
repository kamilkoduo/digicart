package config

import (
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
)

/* vars */
var RedisClient *redis.Client

/* app consts*/
const defaultAppAddress = "0.0.0.0:8080"

/* Redis consts */
const defaultRedisAddress = "0.0.0.0:6379"
const defaultRedisPassword = ""
const defaultRedisDB = 0

var AppAddress = func() string {
	val, found := os.LookupEnv("APP_ADDRESS")
	if !found {
		val = defaultAppAddress
	}
	return val
}()

var RedisAddress = func() string {
	val, found := os.LookupEnv("REDIS_ADDRESS")
	if !found {
		val = defaultRedisAddress
	}
	return val
}()

var RedisPassword = func() string {
	val, found := os.LookupEnv("REDIS_PASSWORD")
	if !found {
		val = defaultRedisPassword
	}
	return val
}()

var RedisDb = func() int {
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
