package api

import (
	"InstantMessaging/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	UploadLocal(c)
}

// UploadLocal 上传文件到本地
func UploadLocal(c *gin.Context) {
	w := c.Writer
	r := c.Request

	// FormFile返回提供表单键的第一个文件
	srcFile, head, err := r.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}

	// 获取上传文件的 Content-Type
	contentType := head.Header.Get("Content-Type")

	// 根据 Content-Type 设置文件后缀
	var suffix string
	switch contentType {
	case "image/jpeg":
		suffix = ".jpg"
	case "image/png":
		suffix = ".png"
	case "audio/mpeg":
		suffix = ".mp3"
	case "video/mp4":
		suffix = ".mp4"
	default:
		suffix = ".dat" // 默认后缀，如果无法识别文件类型
	}

	oFileName := head.Filename           // 获取上传文件的原始文件名
	tem := strings.Split(oFileName, ".") // 使用"."将原始文件名分割成多个部分
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1] // 如果原始文件名包含"."，则获取最后一个"."后面的部分作为文件后缀
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix) // 生成新的文件名，格式为：当前时间戳 + 四位随机数 + 文件后缀
	dstFile, err := os.Create("../front/asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {
			log.Printf("failed, err: %v\n", err)
		}
	}(dstFile)

	// 将源文件内容复制到目标文件中
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	url := "./asset/upload/" + fileName
	utils.RespOK(w, url, "上传文件成功")
}
