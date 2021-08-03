package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/variou/variou-gin/ioc"
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
func (this *Variou) Launch() {
	var port int32 = 8080
	if config := Injector.BeanFactory.Get((*SysConfig)(nil)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	this.applyAll()
	getCronTask().Start()
	this.Run(fmt.Sprintf(":%d", port))
}
func (this *Variou) Handle(httpMethod, relativePath string, handler interface{}) *Variou {
	if h := Convert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, method := range methods {
			getInnerRouter().addRoute(method, this.getPath(relativePath), h) // for future
			this.g.Handle(method, relativePath, h)
		}

	}
	return this
}
func (this *Variou) getPath(relativePath string) string {
	g := "/" + this.currentGroup
	if g == "/" {
		g = ""
	}
	g = g + relativePath
	g = strings.Replace(g, "//", "/", -1)
	return g
}
func (this *Variou) HandleWithFairing(httpMethod, relativePath string, handler interface{}, fairings ...Fairing) *Variou {
	if h := Convert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, f := range fairings {
			Injector.BeanFactory.Apply(f) // set IoC appyly for fairings--- add by shenyi 2020-6-17
		}
		for _, method := range methods {
			getInnerRouter().addRoute(method, this.getPath(relativePath), fairings) //for future
			this.g.Handle(method, relativePath, h)
		}

	}
	return this
}

// 注册中间件
func (this *Variou) Attach(f ...Fairing) *Variou {
	for _, f1 := range f {
		Injector.BeanFactory.Set(f1)
	}
	getFairingHandler().AddFairing(f...)
	return this
}

func (this *Variou) Beans(beans ...Bean) *Variou {
	for _, bean := range beans {
		this.exprData[bean.Name()] = bean
		Injector.BeanFactory.Set(bean)
	}
	return this
}
func (this *Variou) Config(cfgs ...interface{}) *Variou {
	Injector.BeanFactory.Config(cfgs...)
	return this
}
func (this *Variou) applyAll() {
	for t, v := range Injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			Injector.BeanFactory.Apply(v.Interface())
		}
	}
}

func (this *Variou) Mount(group string, classes ...IClass) *Variou {
	this.g = this.Group(group)
	for _, class := range classes {
		this.currentGroup = group
		class.Build(this)
		//this.beanFactory.inject(class)
		this.Beans(class)
	}
	return this
}

//0/3 * * * * *  //增加定时任务
func (this *Variou) Task(cron string, expr interface{}) *Variou {
	var err error
	if f, ok := expr.(func()); ok {
		_, err = getCronTask().AddFunc(cron, f)
	} else if exp, ok := expr.(Expr); ok {
		_, err = getCronTask().AddFunc(cron, func() {
			_, expErr := ExecExpr(exp, this.exprData)
			if expErr != nil {
				log.Println(expErr)
			}
		})
	}

	if err != nil {
		log.Println(err)
	}
	return this
}
