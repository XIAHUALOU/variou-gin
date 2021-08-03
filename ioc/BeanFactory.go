package Injector

import (
	"github.com/XIAHUALOU/variou-gin/expr/src/expr"
	"reflect"
)

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactory()
}

type BeanFactoryImpl struct {
	beanMapper BeanMapper
	ExprMap    map[string]interface{}
}

func (self *BeanFactoryImpl) Set(vlist ...interface{}) {
	if vlist == nil || len(vlist) == 0 {
		return
	}
	for _, v := range vlist {
		self.beanMapper.add(v)
	}
}
func (self *BeanFactoryImpl) Config(cfgs ...interface{}) {
	for _, cfg := range cfgs {
		t := reflect.TypeOf(cfg)
		if t.Kind() != reflect.Ptr {
			panic("required ptr object") //必须是指针对象
		}
		if t.Elem().Kind() != reflect.Struct {
			continue
		}
		self.Set(cfg)                       //把config本身也加入bean
		self.ExprMap[t.Elem().Name()] = cfg //自动构建 ExprMap
		self.Apply(cfg)                     //处理依赖注入 (new)
		v := reflect.ValueOf(cfg)
		for i := 0; i < t.NumMethod(); i++ {
			method := v.Method(i)
			callRet := method.Call(nil)

			if callRet != nil && len(callRet) == 1 {
				self.Set(callRet[0].Interface())
			}
		}
	}
}
func (self *BeanFactoryImpl) Get(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	get_v := self.beanMapper.get(v)
	if get_v.IsValid() {
		return get_v.Interface()
	}
	return nil
}
func (self *BeanFactoryImpl) GetBeanMapper() BeanMapper {
	return self.beanMapper
}

//处理依赖注入
func (self *BeanFactoryImpl) Apply(bean interface{}) {
	if bean == nil {
		return
	}
	v := reflect.ValueOf(bean)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if v.Field(i).CanSet() && field.Tag.Get("inject") != "" {
			if field.Tag.Get("inject") != "-" { //多例模式
				ret := expr.BeanExpr(field.Tag.Get("inject"), self.ExprMap)
				if ret != nil && !ret.IsEmpty() {
					retValue := ret[0]
					if retValue != nil {
						v.Field(i).Set(reflect.ValueOf(retValue))
						self.Apply(retValue)
					}
				}
			} else { //单例模式
				if get_v := self.Get(field.Type); get_v != nil {
					v.Field(i).Set(reflect.ValueOf(get_v))
					self.Apply(get_v)
				}
			}
		}
	}
}
func NewBeanFactory() *BeanFactoryImpl {
	return &BeanFactoryImpl{beanMapper: make(BeanMapper), ExprMap: make(map[string]interface{})}
}
