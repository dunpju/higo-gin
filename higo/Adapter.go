package higo

import (
	"fmt"
	"github.com/dengpju/higo-config/config"
	"github.com/dengpju/higo-logger/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"sync"
)

var (
	orm      *Gorm
	onceGorm sync.Once
)

type Gorm struct {
	*gorm.DB
}

func NewGorm() *Gorm {
	onceGorm.Do(func() {
		confDefault := config.Db("DB.DEFAULT").(*config.Configure)
		args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			confDefault.Get("USERNAME").(string),
			confDefault.Get("PASSWORD").(string),
			confDefault.Get("HOST").(string),
			confDefault.Get("PORT").(string),
			confDefault.Get("DATABASE").(string),
			confDefault.Get("CHARSET").(string),
		)
		db, err := gorm.Open(confDefault.Get("DRIVER").(string), args)
		if err != nil {
			log.Fatal(err)
		}
		logger.Logrus.Infoln(fmt.Sprintf("DB %s:%s Connection success!", confDefault.Get("HOST").(string),
			confDefault.Get("PORT").(string)))
		db.SingularTable(true)
		db.DB().SetMaxIdleConns(5)
		db.DB().SetMaxOpenConns(10)
		orm = &Gorm{DB: db}
	})
	return orm
}
