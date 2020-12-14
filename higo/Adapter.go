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
	_db := Config("DB")
	confDefault := _db.Configure("DEFAULT")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		confDefault.Value("USERNAME"),
		confDefault.Value("PASSWORD"),
		confDefault.Value("HOST"),
		confDefault.Value("PORT"),
		confDefault.Value("DATABASE"),
		confDefault.Value("CHARSET"),
		)
	db, err := gorm.Open(confDefault.Value("DRIVER"), args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("连接数据库成功")
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)
	return &Gorm{DB: db}
}


