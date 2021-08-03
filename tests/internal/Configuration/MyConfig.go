package Configuration

import (
	"github.com/jinzhu/gorm"
	"github.com/variou/variou-gin/tests/internal/Services"
	"log"
	"time"
)

type MyConfig struct {
}

func NewMyConfig() *MyConfig {
	return &MyConfig{}
}
func (this *MyConfig) Test() *Services.TestService {
	return Services.NewTestService("mytest")
}
func (this *MyConfig) Naming() *Services.NameService {
	return Services.NewNameService("variou")
}
func (this *MyConfig) GormDB() *gorm.DB {
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
