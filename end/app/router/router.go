package router

import (
	"InstantMessaging/app/api"
	g "InstantMessaging/app/global"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	//静态资源
	r.Static("/asset", "../front/asset/")
	r.StaticFile("/cqupt.ico", "../front/asset/images/cqupt.ico")
	r.LoadHTMLGlob("../front/views/**/*")

	r.GET("/", api.GetIndex)
	r.GET("/toRegister", api.ToRegister)
	r.GET("/toChat", api.ToChat)
	r.GET("/chat", api.Chat)

	UserRouter := r.Group("/user")
	{
		UserRouter.GET("/getUserList", api.GetUserList)
		UserRouter.POST("/register", api.Register)
		UserRouter.POST("/login", api.Login)
		UserRouter.DELETE("/deleteUser", api.DeleteUser)
		UserRouter.POST("/updateUser", api.UpdateUser)
		UserRouter.POST("/findUser", api.FindUser)
	}

	ContactRouter := r.Group("/contact")
	{
		//ContactRouter.Use(middleware.JWTAuthMiddleware())
		ContactRouter.POST("/addFriend", api.AddFriend)         //加好友
		ContactRouter.POST("/loadFriend", api.LoadFriend)       //好友列表
		ContactRouter.DELETE("/deleteFriend", api.DeleteFriend) //删好友

		ContactRouter.POST("/createCommunity", api.CreateCommunity)   //建群
		ContactRouter.POST("/loadCommunity", api.LoadCommunity)       //群列表
		ContactRouter.POST("/joinCommunity", api.JoinCommunity)       //加群
		ContactRouter.PUT("/updateCommunity", api.UpdateCommunity)    //更改群信息
		ContactRouter.DELETE("/leaveCommunity", api.LeaveCommunity)   //退群
		ContactRouter.DELETE("/deleteCommunity", api.DeleteCommunity) //解散群

		ContactRouter.POST("/redisMsg", api.RedisMsg) //缓存信息
	}

	r.POST("/attach/upload", api.Upload) //上传文件

	g.Logger.Infof("initialize routers successfully")
	return r

}
