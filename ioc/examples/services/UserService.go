package services

import (
	"fmt"
)

type UserService struct {
	Order *OrderService `inject:"-"`
}

func NewUserService() *UserService {
	return &UserService{}
}
func (self *UserService) GetUserInfo(uid int) {
	fmt.Println("GetUserInfo")

}
func (self *UserService) GetOrderInfo(uid int) {
	//self.Order.GetOrderInfo(uid)
	fmt.Println("获取用户ID=", uid, "的订单信息")
}
