package higo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type Gorm struct {
	*gorm.DB
}

func NewGorm() *Gorm {
	config := Container().Config("DB")
	conf := config["DEFAULT"].(map[interface {}]interface{})
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf["USERNAME"].(string),
		conf["PASSWORD"].(string),
		conf["HOST"].(string),
		conf["PORT"].(string),
		conf["DATABASE"].(string),
		conf["CHARSET"].(string),
		)
	db, err := gorm.Open(conf["DRIVER"].(string), args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接数据库成功")
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)
	return &Gorm{DB: db}
}


