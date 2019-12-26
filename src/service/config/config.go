package config

import (
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
)

// RedisClient ...
var RedisClient = func() *redis.Client {
	log.Printf("\nStarting services. Expecting Redis on: %v\n\n", redisAddress)
	rc := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       redisDb,
	})
	_, err := rc.Ping().Result()
	if err != nil {
		log.Fatal("Unable to connect")
	} else {
		log.Printf("Successfully pinged Redis. Ready to accept connections.\n")
	}
	return rc
}()


/* Redis consts */
const defaultRedisAddress = "0.0.0.0:6379"
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
