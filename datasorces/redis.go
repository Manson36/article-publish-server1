package datasorces

import (
	"errors"
	"github.com/article-publish-server1/config"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type Rds struct {
	*redis.Pool
}

func (r *Rds) initDB() {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis.Host+":"+config.Redis.Port)
			if err != nil {
				log.Fatal("redis client dial fail, errmsg:", err.Error())
				return nil, err
			}

			if config.Redis.Password != "" {
				if _, err := c.Do("AUTH", config.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", config.Redis.DB); err != nil {
				c.Close()
				return nil, err
			}

			return c, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000,
		MaxActive:       10000,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	r.Pool = &pool
}

func (r *Rds) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.Get()
	defer func() {
		if e := conn.Close(); e != nil {
		}
	}()

	reply, err = conn.Do(commandName, args...)
	if err != nil {
		if e := conn.Err(); e != nil {
			return nil, errors.New("redis do error, conn errmsg:" + e.Error() + "op errmsg:" + err.Error())
		}
	}

	return reply, err
}

var RDS = &Rds{}
