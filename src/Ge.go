package src

import (
	"github.com/XIAHUALOU/variou-gin/middleWares"
	"github.com/gin-gonic/gin"
)

//The core object of the whole scaffold
type Variou struct { // The nae is the abbreviation of Gin Easy, which means it is more convenient to use gin
	*gin.Engine
	group *gin.RouterGroup
	*Dependency
}

//init Variou and set global error handler middleware
func StartGE() *Variou {
	g := &Variou{Engine: gin.New(), Dependency: NewDependency()}
	g.Use(ErrorHandler())
	return g
}

//start server method
func (self *Variou) Launch() {
	self.Run(SERVER_ADDRESS)
}

//This method is the core of the scaffold. It is mainly for the convenience of returning any type of business results,
// and the binding of groups is done here
func (self *Variou) Handle(httpMethod, relativePath string, handler interface{}) *Variou {
	if h := Convert(handler); h != nil {
		self.group.Handle(httpMethod, relativePath, h)
	}
	return self
}

//Variou's middleware
func (self *Variou) AddMid(middleWares ...middleWares.Mid) *Variou {
	for _, mid := range middleWares {
		mid := mid
		self.Use(func(context *gin.Context) {
			err := mid.BeforeRequest(context)
			if err != nil {
				context.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			} else {
				context.Next()
			}
		})
	}
	return self
}

//This method is mainly used to add routes and inject dependencies into controllers
func (self *Variou) AddController(group string, controllers ...Controller) *Variou {
	self.group = self.Group(group)
	for _, controller := range controllers {
		controller.Build(self)
		self.BindDep(controller)
	}
	return self
}

func (self *Variou) PrepareDeps(deps ...interface{}) *Variou {
	self.SetDep(deps...)
	return self
}
