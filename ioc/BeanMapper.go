package Injector

import (
	"reflect"
)

type BeanMapper map[reflect.Type]reflect.Value

func (self BeanMapper) add(bean interface{}) {
	t := reflect.TypeOf(bean)
	if t.Kind() != reflect.Ptr {
		panic("require ptr object")
	}
	self[t] = reflect.ValueOf(bean)
}
func (self BeanMapper) get(bean interface{}) reflect.Value {
	var t reflect.Type
	if bt, ok := bean.(reflect.Type); ok {
		t = bt
	} else {
		t = reflect.TypeOf(bean)
	}
	if v, ok := self[t]; ok {
		return v
	}
	//处理接口 继承
	for k, v := range self {
		if t.Kind() == reflect.Interface && k.Implements(t) {
			return v
		}
	}

	return reflect.Value{}
}
