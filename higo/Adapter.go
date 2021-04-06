package higo

import (
	"fmt"
	"gitee.com/dengpju/higo-configure/configure"
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
		_db := configure.Config("DB")
		confDefault := _db.Configure("DEFAULT")
		args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			confDefault.String("USERNAME"),
			confDefault.String("PASSWORD"),
			confDefault.String("HOST"),
			confDefault.String("PORT"),
			confDefault.String("DATABASE"),
			confDefault.String("CHARSET"),
		)
		db, err := gorm.Open(confDefault.String("DRIVER"), args)
		if err != nil {
			log.Fatal(err)
		}
		logger.Logrus.Infoln(fmt.Sprintf("DB %s:%s Connection success!", confDefault.String("HOST"),
			confDefault.String("PORT")))
		db.SingularTable(true)
		db.DB().SetMaxIdleConns(5)
		db.DB().SetMaxOpenConns(10)
		orm = &Gorm{DB: db}
	})
	return orm
}
