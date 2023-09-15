package api

import (
	"InstantMessaging/app/internal/model"
	"InstantMessaging/app/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"text/template"
)

func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("../front/index.html", "../front/views/chat/head.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		return
	}
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("../front/views/user/register.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "register")
	if err != nil {
		return
	}
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("../front/views/chat/index.html",
		"../front/views/chat/head.html",
		"../front/views/chat/foot.html",
		"../front/views/chat/tabmenu.html",
		"../front/views/chat/concat.html",
		"../front/views/chat/group.html",
		"../front/views/chat/profile.html",
		"../front/views/chat/createcom.html",
		"../front/views/chat/userinfo.html",
		"../front/views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")

	// 创建一个User对象，并设置其ID和Identity属性
	u := model.User{}
	u.ID = uint(userId)
	u.Identity = token

	// 使用解析后的模板执行HTML渲染，并将User对象传递给模板
	err = ind.Execute(c.Writer, u)
	if err != nil {
		return
	}
}

func Chat(c *gin.Context) {
	service.Chat(c.Writer, c.Request)
}
