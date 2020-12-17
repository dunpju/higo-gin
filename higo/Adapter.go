package higo

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

type Gorm struct {
	*gorm.DB
}

func NewGorm() *Gorm {
	_db := Config("DB")
	confDefault := _db.Configure("DEFAULT")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		confDefault.StrValue("USERNAME"),
		confDefault.StrValue("PASSWORD"),
		confDefault.StrValue("HOST"),
		confDefault.StrValue("PORT"),
		confDefault.StrValue("DATABASE"),
		confDefault.StrValue("CHARSET"),
		)
	db, err := gorm.Open(confDefault.StrValue("DRIVER"), args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接数据库成功")
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)
	return &Gorm{DB: db}
}

var RedisPool *redis.Pool

type Redis struct {
	*redis.Pool
}

func NewRedisPool() *redis.Pool {
	_redis := Config("REDIS")
	confDefault := _redis.Configure("DEFAULT")
	pool := confDefault.Configure("POOL")
	return &redis.Pool{
		MaxIdle: pool.IntValue("MAX_IDLE"),
		IdleTimeout: time.Duration(pool.IntValue("MAX_IDLE_TIME")) * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp",
				fmt.Sprintf("%s:%s",confDefault.StrValue("HOST"),confDefault.StrValue("PORT")),
				redis.DialDatabase(confDefault.IntValue("DB")),
				redis.DialPassword(confDefault.StrValue("AUTH")),
				)
		},
	}
} 




