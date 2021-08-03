package classes

import (
	"fmt"
	"github.com/XIAHUALOU/variou-gin/core"
	"github.com/XIAHUALOU/variou-gin/tests/internal/Services"
	"github.com/XIAHUALOU/variou-gin/tests/internal/fairing"
	"github.com/gin-gonic/gin"
)

type MyError struct {
	Code    int
	Message string
}

func NewMyError(code int, message string) *MyError {
	return &MyError{Code: code, Message: message}
}
func (*MyError) Name() string {
	return "myerror"
}

type IndexClass struct {
	MyTest  *Services.TestService `inject:"-"`
	MyTest2 *Services.TestService
	Age     *core.Value `prefix:"user.age"`
}

func NewIndexClass() *IndexClass {

	return &IndexClass{}
}
func (self *IndexClass) GetIndex(ctx *gin.Context) string {
	self.MyTest.Naming.ShowName()
	return "IndexClass"
}
func (self *IndexClass) TestA(c *gin.Context) core.Json {

	return gin.H{"message": "testa"}
}
func (self *IndexClass) Test(ctx *gin.Context) core.Json {
	//fmt.Println("name is", ctx.PostForm("name"))

	//ctx.Set(core.HTTP_STATUS, 503)
	panic(NewMyError(1800, "oh shit"))
	//fmt.Println(self.Age.String())
	return NewDataModel(101, "wfew")
}
func (self *IndexClass) TestUsers(ctx *gin.Context) core.Query {

	return core.SimpleQuery("select * from users").WithMapping(map[string]string{
		"user_name": "uname",
	}).WithKey("result")
}
func (self *IndexClass) TestUserDetail(ctx *gin.Context) core.Json {
	ret := core.SimpleQuery("select * from users where user_id=?").
		WithArgs(ctx.Param("id")).WithMapping(map[string]string{
		"usr": "user",
	}).WithFirst().WithKey("result").Get()

	fmt.Printf("%T", ret.(gin.H)["result"].(map[string]interface{}))
	return ret
}
func (self *IndexClass) IndexVoid(c *gin.Context) (void core.Void) {
	c.JSON(200, gin.H{"message": "void"})
	return
}
func (self *IndexClass) Build(goft *core.Variou) {
	goft.HandleWithFairing("GET", "/",
		self.GetIndex, fairing.NewIndexFairing()).
		Handle("GET", "/users", self.TestUsers).
		Handle("GET", "/users/:id", self.TestUserDetail).
		Handle("GET", "/test", self.Test)
}
func (self *IndexClass) Name() string {
	return "IndexClass"
}
