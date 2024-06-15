package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-img/conf"
	"log"
	"time"
)

var rdb *redisCache

type redisCache struct {
	Enable            bool
	Client            *redis.Client
	Context           context.Context
	DefaultExpireTime time.Duration
}

func init() {
	rdb = new(redisCache)
	rdb.Enable = conf.AppConf.CacheConf.Enable
	rdb.Context = context.Background()
	rdb.DefaultExpireTime = time.Duration(conf.AppConf.CacheConf.ExpireTime) * time.Second
	rdb.Client = nil
	if rdb.Enable {
		rdb.Client = redis.NewClient(&redis.Options{
			Addr:     conf.AppConf.CacheConf.Address,
			Password: conf.AppConf.CacheConf.Password,
			DB:       0,
		})
		if _, err := rdb.Client.Ping(rdb.Context).Result(); err != nil {
			log.Fatalln("connect redis server failed:", err)
		}
	}
}

func NewRedisCache() *redisCache {
	return rdb
}

func (rc *redisCache) Set(key string, value interface{}) error {
	if uint32(len(value.([]byte))) <= conf.AppConf.CacheConf.MaxCache {
		return rc.Client.Set(rc.Context, key, value, rc.DefaultExpireTime).Err()
	}
	return fmt.Errorf("greater than max cahce size")
}

func (rc *redisCache) Get(key string) ([]byte, error) {
	val, err := rc.Client.Get(rc.Context, key).Bytes()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	return val, nil
}

func (rc *redisCache) Del(key string) error {
	return rc.Client.Del(rc.Context, key).Err()
}
