package fairing

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type IndexFairing struct{}

func NewIndexFairing() *IndexFairing {
	return &IndexFairing{}
}

func (self *IndexFairing) OnRequest(ctx *gin.Context) error {
	if v, exists := ctx.Get("name"); exists {
		v = fmt.Sprintf("%v,self is index name")
		ctx.Set("name", v)
	}
	return nil
}
func (self *IndexFairing) OnResponse(ret interface{}) (interface{}, error) {
	if str, ok := ret.(string); ok {
		str = str + "_index"
		return str, nil
	}

	return ret, nil
}
