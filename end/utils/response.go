package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type H struct {
	Code  int
	Msg   string
	Data  interface{}
	Rows  interface{}
	Total interface{}
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, http.StatusInternalServerError, nil, msg)
}

func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, http.StatusOK, data, msg)
}

func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, http.StatusOK, data, total)
}

func RespSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
	})
}

func RespFailed(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": message,
	})
}
