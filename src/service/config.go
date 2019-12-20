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

/* app consts*/
const defaultAppAdress = "0.0.0.0:8080"

/* Redis consts */
const defaultRedisAddress = "0.0.0.0:6379"
const defaultRedisPassword = ""
const defaultRedisDB = 0

var AppAddress = func() string {
	val, found := os.LookupEnv("APP_ADDRESS")
	if !found {
		val = defaultAppAdress
	}
	return val
}()

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
