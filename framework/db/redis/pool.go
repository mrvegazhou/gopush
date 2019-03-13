package redigoUtil

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopush/conf"
	"gopush/framework/log"
	"time"
)

var pool *redis.Pool

func initPool(config *conf.MainConfig) {
	pool = &redis.Pool{
		MaxIdle:     config.Redis.MaxIdle,
		MaxActive:   config.Redis.MaxActive,
		IdleTimeout: 30 * time.Second,
		Wait:        config.Redis.Wait,
		Dial: func() (redis.Conn, error) {
			return setDialog(config)
		},
	}
}

func setDialog(config *conf.MainConfig) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Redis.Host, config.Port))
	if err != nil {
		logger.Sugar.Error(fmt.Sprintf("init redis failed! %v", config))
	}
	if len(config.Redis.Password) != 0 {
		if _, err := conn.Do("AUTH", config.Redis.Password); err != nil {
			conn.Close()
			logger.Sugar.Error(err)
		}
	}
	if _, err := conn.Do("SELECT", config.Redis.DbNum); err != nil {
		conn.Close()
		logger.Sugar.Error(err)
	}
	r, err := redis.String(conn.Do("PING"))
	if err != nil || r != "PONG" {
		panic("连接失败")
	}

	return conn, nil
}
