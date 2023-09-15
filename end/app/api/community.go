package api

import (
	"InstantMessaging/app/internal/model"
	"InstantMessaging/app/internal/service"
	"InstantMessaging/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateCommunity 新建群
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))
	name := c.PostForm("name")
	icon := c.PostForm("icon")
	cate := c.PostForm("cate")
	desc := c.PostForm("desc")
	community := model.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Type = cate
	community.Desc = desc
	ok, msg := service.CreateCommunity(community)
	if ok {
		utils.RespOK(c.Writer, http.StatusOK, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// UpdateCommunity 更新群信息
func UpdateCommunity(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	comId := c.PostForm("comId")
	name := c.PostForm("name")
	icon := c.PostForm("icon")
	cate := c.PostForm("cate")
	desc := c.PostForm("desc")

	community := model.Community{}
	community.OwnerId = uint(userId)
	community.Name = name
	community.Img = icon
	community.Type = cate
	community.Desc = desc

	ok, msg := service.UpdateCommunity(community, comId)
	if ok {
		utils.RespOK(c.Writer, http.StatusOK, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// DeleteCommunity 删除群
func DeleteCommunity(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	comId := c.PostForm("comId")

	ok, msg := service.DeleteCommunity(uint(userId), comId)
	if ok {
		utils.RespOK(c.Writer, http.StatusOK, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// LoadCommunity 加载群列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))
	//	name := c.PostForm("name")
	communities := service.LoadCommunity(uint(ownerId))
	utils.RespOKList(c.Writer, communities, len(communities))
}

// JoinCommunity 加入群 userId uint, comId uint
func JoinCommunity(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	comId := c.PostForm("comId")
	//	name := c.PostForm("name")

	ok, msg := service.JoinCommunity(uint(userId), comId)
	if ok {
		utils.RespOK(c.Writer, http.StatusOK, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// LeaveCommunity 退群
func LeaveCommunity(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	comId := c.PostForm("comId")

	ok, msg := service.LeaveCommunity(uint(userId), comId)
	if ok {
		utils.RespOK(c.Writer, http.StatusOK, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}
