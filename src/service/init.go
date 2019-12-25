package service

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/kamilkoduo/digicart/src/service/config"
	"log"
	"strings"
)

func init() {
	log.Println("Starting Cart Service")
	log.Println("Active variables:")
	log.Println("  REDIS_ADDRESS     = ", config.RedisAddress)
	log.Println("  REDIS_PASSWORD    = ", strings.Repeat("*", len(config.RedisPassword)))
	log.Println("  REDIS_DB          = ", config.RedisDb)
	err := pingRedis()
	if err != nil {
		log.Fatalln("Unable to connect to the database at " + config.RedisAddress + ": " + err.Error())
	}

	log.Println("Connection to the database was established")
}

func pingRedis() (err error) {
	config.RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	})
	_, err = config.RedisClient.Ping().Result()
	if err != nil {
		return errors.New("Unable to ping Redis database: " + err.Error())
	}
	return nil
}
