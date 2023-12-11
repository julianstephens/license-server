package redis

import (
	"context"
	"time"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var STORE *persist.RedisStore

func Init() error {
	db := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	err := db.Set(context.TODO(), "test", "hello world", 2*time.Minute).Err()
	if err != nil {
		return err
	}

	store := persist.NewRedisStore(db)

	RDB = db
	STORE = store
	return nil
}

func GetStore() *persist.RedisStore {
	return STORE
}
