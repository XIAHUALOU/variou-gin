package main

import (
	"fmt"
	"github.com/variou/variou-gin/ioc"
	"github.com/variou/variou-gin/ioc/examples/Config"
	"github.com/variou/variou-gin/ioc/examples/services"
)

func main() {
	serviceConfig := Config.NewServiceConfig()

	Injector.BeanFactory.Config(serviceConfig) //展开方法
	//  BeanFactory.Set()
	{
		//这里 测试 userServices
		userService := services.NewUserService()
		Injector.BeanFactory.Apply(userService) //处理依赖
		fmt.Println(userService.Order.Name())
		userService.GetUserInfo(3)

	}
}
