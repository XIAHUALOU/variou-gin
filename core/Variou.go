package core

import (
	"fmt"
	"github.com/XIAHUALOU/variou-gin/ioc"
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
	"strings"
	"sync"
)

type Bean interface {
	Name() string
}

var Empty = &struct{}{}
var innerRouter *GoftTree // inner tree node . backup httpmethod and path
var innerRouter_once sync.Once

func getInnerRouter() *GoftTree {
	innerRouter_once.Do(func() {
		innerRouter = NewGoftTree()
	})
	return innerRouter
}

type Variou struct {
	*gin.Engine
	g            *gin.RouterGroup // 保存 group对象
	exprData     map[string]interface{}
	currentGroup string // temp-var for group string
}

func Ignite(ginMiddlewares ...gin.HandlerFunc) *Variou {
	g := &Variou{Engine: gin.New(),
		exprData: map[string]interface{}{},
	}
	g.Use(ErrorHandler()) //强迫加载的异常处理中间件
	for _, handler := range ginMiddlewares {
		g.Use(handler)
	}
	config := InitConfig()
	Injector.BeanFactory.Set(g)      // inject self
	Injector.BeanFactory.Set(config) // add global into (new)BeanFactory
	Injector.BeanFactory.Set(NewGPAUtil())
	if config.Server.Html != "" {
		g.LoadHTMLGlob(config.Server.Html)
	}
	return g
}
func (self *Variou) Launch() {
	var port int32 = 8080
	if config := Injector.BeanFactory.Get((*SysConfig)(nil)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	self.applyAll()
	getCronTask().Start()
	self.Run(fmt.Sprintf(":%d", port))
}
func (self *Variou) Handle(httpMethod, relativePath string, handler interface{}) *Variou {
	if h := Convert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, method := range methods {
			getInnerRouter().addRoute(method, self.getPath(relativePath), h) // for future
			self.g.Handle(method, relativePath, h)
		}

	}
	return self
}
func (self *Variou) getPath(relativePath string) string {
	g := "/" + self.currentGroup
	if g == "/" {
		g = ""
	}
	g = g + relativePath
	g = strings.Replace(g, "//", "/", -1)
	return g
}
func (self *Variou) HandleWithFairing(httpMethod, relativePath string, handler interface{}, fairings ...Fairing) *Variou {
	if h := Convert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, f := range fairings {
			Injector.BeanFactory.Apply(f) // set IoC appyly for fairings--- add by shenyi 2020-6-17
		}
		for _, method := range methods {
			getInnerRouter().addRoute(method, self.getPath(relativePath), fairings) //for future
			self.g.Handle(method, relativePath, h)
		}

	}
	return self
}

// 注册中间件
func (self *Variou) Attach(f ...Fairing) *Variou {
	for _, f1 := range f {
		Injector.BeanFactory.Set(f1)
	}
	getFairingHandler().AddFairing(f...)
	return self
}

func (self *Variou) Beans(beans ...Bean) *Variou {
	for _, bean := range beans {
		self.exprData[bean.Name()] = bean
		Injector.BeanFactory.Set(bean)
	}
	return self
}
func (self *Variou) Config(cfgs ...interface{}) *Variou {
	Injector.BeanFactory.Config(cfgs...)
	return self
}
func (self *Variou) applyAll() {
	for t, v := range Injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			Injector.BeanFactory.Apply(v.Interface())
		}
	}
}

func (self *Variou) Mount(group string, classes ...IClass) *Variou {
	self.g = self.Group(group)
	for _, class := range classes {
		self.currentGroup = group
		class.Build(self)
		//self.beanFactory.inject(class)
		self.Beans(class)
	}
	return self
}

//0/3 * * * * *  //增加定时任务
func (self *Variou) Task(cron string, expr interface{}) *Variou {
	var err error
	if f, ok := expr.(func()); ok {
		_, err = getCronTask().AddFunc(cron, f)
	} else if exp, ok := expr.(Expr); ok {
		_, err = getCronTask().AddFunc(cron, func() {
			_, expErr := ExecExpr(exp, self.exprData)
			if expErr != nil {
				log.Println(expErr)
			}
		})
	}

	if err != nil {
		log.Println(err)
	}
	return self
}
