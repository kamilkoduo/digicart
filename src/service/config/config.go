package config

import (
	"github.com/kamilkoduo/digicart/src/service/db/api"
	"github.com/kamilkoduo/digicart/src/service/db/service"
)

// CartDBAPIServer ...
var CartDBAPIServer api.CartDBAPI = service.CartDBRedisServer{}
