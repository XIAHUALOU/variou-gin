package Configuration

import (
	"github.com/XIAHUALOU/variou-gin/core"
	"github.com/XIAHUALOU/variou-gin/tests/internal/classes"
)

type RouterConfig struct {
	Goft       *core.Variou        `inject:"-"`
	IndexClass *classes.IndexClass `inject:"-"`
}

func NewRouterConfig() *RouterConfig {
	return &RouterConfig{}
}
func (self *RouterConfig) IndexRoutes() interface{} {
	self.Goft.Handle("GET", "/a", self.IndexClass.TestA)
	self.Goft.Handle("GET", "/b", self.IndexClass.TestA)
	self.Goft.Handle("GET", "/void", self.IndexClass.IndexVoid)
	return core.Empty
}
