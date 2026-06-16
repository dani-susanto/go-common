package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func New(
	username string,
	password string,
	host string,
	port string,
	database string,
) (*redis.Client, error) {
	dsn, err := redis.ParseURL(
		fmt.Sprintf(
			"redis://%s:%s@%s:%s/%s",
			username,
			password,
			host,
			port,
			database,
		),
	)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(dsn), nil
}
