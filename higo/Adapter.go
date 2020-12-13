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
	config := Config("DB").(map[string]interface{})
	conf := config["DEFAULT"].(map[string]interface{})
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf["USERNAME"],
		conf["PASSWORD"],
		conf["HOST"],
		conf["PORT"],
		conf["DATABASE"],
		conf["CHARSET"],
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


