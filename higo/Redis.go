package higo

import (
	"github.com/dengpju/higo-config/config"
	hiredis "github.com/dengpju/higo-redis/redis"
	"github.com/gomodule/redigo/redis"
	"sync"
)

var RedisPool *redis.Pool

var redisOnce sync.Once

func InitRedisPool() *redis.Pool {
	redisOnce.Do(func() {
		confDefault := config.Db("Redis.Default").(*config.Configure)
		pool := confDefault.Get("Pool").(*config.Configure)
		RedisPool = hiredis.New(
			hiredis.NewPoolConfigure(
				hiredis.PoolHost(confDefault.Get("Host").(string)),
				hiredis.PoolPort(confDefault.Get("Port").(int)),
				hiredis.PoolAuth(confDefault.Get("Auth").(string)),
				hiredis.PoolDb(confDefault.Get("Db").(int)),
				hiredis.PoolMaxConnections(pool.Get("Max_Connections").(int)),
				hiredis.PoolMaxIdle(pool.Get("Max_Idle").(int)),
				hiredis.PoolMaxIdleTime(pool.Get("Max_Idle_Time").(int)),
				hiredis.PoolMaxConnLifetime(pool.Get("Max_Conn_Lifetime").(int)),
				hiredis.PoolWait(pool.Get("Wait").(bool)),
			))
		Redis = hiredis.Redis
	})
	return RedisPool
}

var Redis hiredis.RedisAdapter

type RedisAdapter struct {
	*hiredis.RedisAdapter
}

func NewRedisAdapter() *RedisAdapter {
	return &RedisAdapter{}
}
