package service

import (
	"errors"
	"github.com/go-redis/redis"
	"log"
	"strings"
)

func init() {
	log.Println("Starting Cart Service")
	log.Println("Active variables:")
	log.Println("  REDIS_ADDRESS     = ", redisAddress)
	log.Println("  REDIS_PASSWORD    = ", strings.Repeat("*", len(redisPassword)))
	log.Println("  REDIS_DB          = ", redisDb)
	err := pingRedis()
	if err != nil {
		log.Fatalln("Unable to connect to the database at " + redisAddress + ": " + err.Error())
	}

	log.Println("Connection to the database was established")
}

func pingRedis() (err error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       redisDb,
	})
	_, err = redisClient.Ping().Result()
	if err != nil {
		return errors.New("Unable to ping Redis database: " + err.Error())
	}
	return nil

}
