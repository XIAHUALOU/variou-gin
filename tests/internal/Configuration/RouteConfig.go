package Configuration

import (
	"github.com/variou/variou-gin/core"
	"github.com/variou/variou-gin/tests/internal/classes"
)

type RouterConfig struct {
	Goft       *core.Variou        `inject:"-"`
	IndexClass *classes.IndexClass `inject:"-"`
}

func NewRouterConfig() *RouterConfig {
	return &RouterConfig{}
}
func (this *RouterConfig) IndexRoutes() interface{} {
	this.Goft.Handle("GET", "/a", this.IndexClass.TestA)
	this.Goft.Handle("GET", "/b", this.IndexClass.TestA)
	this.Goft.Handle("GET", "/void", this.IndexClass.IndexVoid)
	return core.Empty
}
