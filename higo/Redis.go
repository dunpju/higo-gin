package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-configure/configure"
	"github.com/dengpju/higo-throw/throw"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var RedisPool *redis.Pool

var redisOnce sync.Once

func InitRedisPool() *redis.Pool {
	redisOnce.Do(func() {
		_redis := configure.Config("REDIS")
		confDefault := _redis.Configure("DEFAULT")
		pool := confDefault.Configure("POOL")
		RedisPool = &redis.Pool {
			MaxActive:   pool.Int("MAX_CONNECTIONS"),
			MaxIdle:     pool.Int("MAX_IDLE"),
			IdleTimeout: time.Duration(pool.Int("MAX_IDLE_TIME")) * time.Second,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial("tcp",
					fmt.Sprintf("%s:%s", confDefault.String("HOST"), confDefault.String("PORT")),
					redis.DialDatabase(confDefault.Int("DB")),
					redis.DialPassword(confDefault.String("AUTH")),
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
		throw.Throw(throw.Message(err), throw.Code(0))
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
		throw.Throw(throw.Message(err), throw.Code(0))
	}
	return true
}