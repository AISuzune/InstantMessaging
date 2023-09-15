package api

import (
	"InstantMessaging/app/internal/model"
	"InstantMessaging/app/internal/service"
	"InstantMessaging/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"
)

func GetUserList(c *gin.Context) {
	data := make([]*model.User, 10)
	data = service.GetUserList()
	utils.RespSuccess(c, data, "用户名已注册！")

}

func Register(c *gin.Context) {

	data := model.User{}
	username := c.PostForm("username")
	password := c.PostForm("password")
	rePassword := c.PostForm("rePassword")
	fmt.Println(username, "  >>>>>>>>>>>  ", password, rePassword)
	// 生成一个随机的盐值
	salt := fmt.Sprintf("%06d", rand.Int31())
	// 检查用户名是否已存在
	u := service.CheckUserIsExist(username)

	if username == "" || password == "" || rePassword == "" {
		utils.RespFailed(c, "用户名或密码不能为空！")
		return
	}
	if u.Username != "" {
		utils.RespFailed(c, "用户名已注册！")
		return
	}
	if password != rePassword {
		utils.RespFailed(c, "两次密码不一致！")
		return
	}

	// 设置用户信息
	u.Username = username
	u.Password = utils.MakePassword(password, salt)
	u.Salt = salt
	fmt.Println(u.Password)
	u.LoginTime = time.Now()
	u.LoginOutTime = time.Now()
	u.HeartbeatTime = time.Now()
	service.CreateUser(u)
	utils.RespSuccess(c, data, "新增用户成功！")

}

func Login(c *gin.Context) {
	data := model.User{}

	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println(username, password)
	u := service.CheckUserIsExist(username)
	if u.Username == "" {
		utils.RespFailed(c, "该用户不存在")
		return
	}

	flag := utils.ValidPassword(password, u.Salt, u.Password)
	if !flag {
		utils.RespFailed(c, "密码不正确")
		return
	}
	pwd := utils.MakePassword(password, u.Salt)
	data = service.CheckUsernameAndPwd(username, pwd)
	utils.RespSuccess(c, data, "登录成功")

}

func DeleteUser(c *gin.Context) {
	u := model.User{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	u.ID = uint(id)
	service.DeleteUser(u)
	utils.RespSuccess(c, u, "删除用户成功！")

}

func UpdateUser(c *gin.Context) {
	u := model.User{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	u.ID = uint(id)
	u.Username = c.PostForm("username")
	u.Password = c.PostForm("password")
	u.Phone = c.PostForm("phone")
	u.Avatar = c.PostForm("icon")
	u.Email = c.PostForm("email")
	fmt.Println("update :", u)

	// 使用 govalidator 对 User 结构体进行验证
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		fmt.Println(err)
		utils.RespFailed(c, "修改参数不匹配！")
	} else {
		service.UpdateUser(u)
		utils.RespSuccess(c, u, "修改用户成功！")
	}

}

func FindUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	data := service.FindUser(uint(userId))
	utils.RespSuccess(c, data, "ok")
}
