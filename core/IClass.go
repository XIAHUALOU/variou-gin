package core

type IClass interface {
	Build(goft *Variou) //参数和方法名必须一致
	Name() string
}
