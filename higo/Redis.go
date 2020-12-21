package higo

import (
	"fmt"
	"github.com/dengpju/higo-throw/throw"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var RedisPool *redis.Pool

var redisOnce sync.Once

func InitRedisPool() *redis.Pool {
	redisOnce.Do(func() {
		_redis := Config("REDIS")
		confDefault := _redis.Configure("DEFAULT")
		pool := confDefault.Configure("POOL")
		RedisPool = &redis.Pool {
			MaxActive:   pool.IntValue("MAX_CONNECTIONS"),
			MaxIdle:     pool.IntValue("MAX_IDLE"),
			IdleTimeout: time.Duration(pool.IntValue("MAX_IDLE_TIME")) * time.Second,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial("tcp",
					fmt.Sprintf("%s:%s", confDefault.StrValue("HOST"), confDefault.StrValue("PORT")),
					redis.DialDatabase(confDefault.IntValue("DB")),
					redis.DialPassword(confDefault.StrValue("AUTH")),
				)
			},
		}
		Redis = RedisAdapter{}
	})
	return RedisPool
}

var Redis RedisAdapter

type RedisAdapter struct {
	conn redis.Conn
}

func NewRedisAdapter() *RedisAdapter {
	return &RedisAdapter{}
}

func (this *RedisAdapter) Connection() *RedisAdapter {
	this.conn = RedisPool.Get()
	return this
}

func (this *RedisAdapter) Set(key string, v interface{}) bool {
	this.Connection()
	defer this.conn.Close()
	_, err := this.conn.Do("set", key, v)
	if err != nil {
		this.conn.Close()
		throw.Throw(err, 0)
	}
	return true
}

func (this *RedisAdapter) Get(key string) (string, error) {
	this.Connection()
	defer this.conn.Close()
	v, err := redis.String(this.conn.Do("get", key))
	return v, err
}

func (this *RedisAdapter) GetByte(key string) ([]byte, error) {
	this.Connection()
	defer this.conn.Close()
	v, err := redis.Bytes(this.conn.Do("get", key))
	return v, err
}

func (this *RedisAdapter) Setex(key string, expire int, data []byte) bool {
	this.Connection()
	defer this.conn.Close()
	_, err := this.conn.Do("setex", key, expire, data)
	if err != nil {
		this.conn.Close()
		throw.Throw(err, 0)
	}
	return true
}