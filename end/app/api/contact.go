package api

import (
	"InstantMessaging/app/internal/service"
	"InstantMessaging/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LoadFriend(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("userId"))
	users := service.LoadFriend(uint(id))
	utils.RespOKList(c.Writer, users, len(users))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	targetName := c.PostForm("targetName")
	//targetId, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	ok, msg := service.AddFriend(uint(userId), targetName)
	if ok {
		utils.RespOK(c.Writer, http.StatusOK, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func DeleteFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	targetName := c.PostForm("targetName")

	ok, msg := service.DeleteFriend(uint(userId), targetName)
	if ok {
		utils.RespOK(c.Writer, http.StatusOK, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := service.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, res, len(res))
}
