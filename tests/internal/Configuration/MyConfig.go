package Configuration

import (
	"github.com/XIAHUALOU/variou-gin/tests/internal/Services"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type MyConfig struct {
}

func NewMyConfig() *MyConfig {
	return &MyConfig{}
}
func (self *MyConfig) Test() *Services.TestService {
	return Services.NewTestService("mytest")
}
func (self *MyConfig) Naming() *Services.NameService {
	return Services.NewNameService("variou")
}
func (self *MyConfig) GormDB() *gorm.DB {
	db, err := gorm.Open("mysql",
		"root:123123@tcp(localhost:3307)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)
	db.DB().SetConnMaxLifetime(time.Second * 30)
	return db
}
