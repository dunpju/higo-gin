package higo

import (
	"fmt"
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-throw/throw"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var RedisPool *redis.Pool

var redisOnce sync.Once

func InitRedisPool() *redis.Pool {
	redisOnce.Do(func() {
		confDefault := config.Get("env.app.REDIS.DEFAULT").(config.Configure)
		pool := confDefault.Get("POOL").(config.Configure)
		RedisPool = &redis.Pool {
			MaxActive:   pool.Get("MAX_CONNECTIONS").(int),
			MaxIdle:     pool.Get("MAX_IDLE").(int),
			IdleTimeout: time.Duration(pool.Get("MAX_IDLE_TIME").(int)) * time.Second,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial("tcp",
					fmt.Sprintf("%s:%s", confDefault.Get("HOST").(string), confDefault.Get("PORT").(string)),
					redis.DialDatabase(confDefault.Get("DB").(int)),
					redis.DialPassword(confDefault.Get("AUTH").(string)),
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