package service

import (
	g "InstantMessaging/app/global"
	"InstantMessaging/app/internal/middleware"
	"InstantMessaging/app/internal/model"
	"InstantMessaging/utils"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/gorm"
	"time"
)

func GetUserList() []*model.User {
	data := make([]*model.User, 0)
	g.MysqlDB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CheckUserIsExist(username string) model.User {
	user := model.User{}
	g.MysqlDB.Where("username = ?", username).First(&user)
	return user
}

func CheckUsernameAndPwd(username string, password string) model.User {
	user := model.User{}
	g.MysqlDB.Where("username = ? and password = ?", username, password).First(&user)

	//token加密
	//str := fmt.Sprintf("%d", time.Now().Unix())
	//temp := utils.MD5Encode(str)
	claim := utils.MyClaims{
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 过期时间
			Issuer:    "ZTY",                                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, _ := token.SignedString(middleware.Secret)
	g.MysqlDB.Model(&user).Where("id = ?", user.ID).Update("identity", tokenString)
	return user
}

func FindUserByPhone(phone string) *gorm.DB {
	user := model.User{}
	return g.MysqlDB.Where("Phone = ?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := model.User{}
	return g.MysqlDB.Where("email = ?", email).First(&user)
}

func CreateUser(user model.User) *gorm.DB {
	return g.MysqlDB.Create(&user)
}

func DeleteUser(user model.User) *gorm.DB {
	return g.MysqlDB.Delete(&user)
}

func UpdateUser(user model.User) *gorm.DB {
	return g.MysqlDB.Model(&user).Updates(model.User{Username: user.Username, Password: user.Password, Phone: user.Phone, Email: user.Email, Avatar: user.Avatar})
}

// FindUser 查找某个用户
func FindUser(id uint) model.User {
	user := model.User{}
	g.MysqlDB.Where("id = ?", id).First(&user)
	return user
}
