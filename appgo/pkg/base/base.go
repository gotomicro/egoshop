package base

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/pkg/code"
)

type RespList struct {
	List       interface{} `json:"list"`
	Pagination struct {
		Current  int `json:"current"`
		PageSize int `json:"pageSize"`
		Total    int `json:"total"`
	} `json:"pagination"`
}

// JSONResult json
type JSONResult struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type JSONResultRaw struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

// JSON 提供了系统标准JSON输出方法。
func JSON(c *gin.Context, Code int, data ...interface{}) {
	result := new(JSONResult)
	result.Code = Code
	info, ok := code.CodeMap[Code]
	if ok {
		result.Message = info
	} else {
		result.Message = "error"
	}

	if len(data) > 0 {
		result.Data = data[0]
	} else {
		result.Data = ""
	}
	c.JSON(http.StatusOK, result)
	return
}

func JSONOK(c *gin.Context, result ...interface{}) {
	j := new(JSONResult)
	j.Code = 0
	j.Message = "ok"
	if len(result) > 0 {
		j.Data = result[0]
	} else {
		j.Data = ""
	}
	c.JSON(http.StatusOK, j)
	return
}

// JSON 提供了系统标准JSON输出方法。
func JSONErr(c *gin.Context, Code int, err error) {
	result := new(JSONResult)
	result.Code = Code
	info, ok := code.CodeMap[Code]
	if ok {
		result.Message = info
	} else {
		result.Message = "error"
	}
	if err != nil {
		fmt.Println("code is", Code, "info is", result.Message, "============== err is", err.Error())
	}

	c.JSON(http.StatusOK, result)
	return
}

func JSONList(c *gin.Context, data interface{}, current, pageSize, total int) {
	j := new(JSONResult)
	j.Code = 0
	j.Message = "ok"
	j.Data = RespList{
		List: data,
		Pagination: struct {
			Current  int `json:"current"`
			PageSize int `json:"pageSize"`
			Total    int `json:"total"`
		}{
			Current:  current,
			PageSize: pageSize,
			Total:    total,
		},
	}
	c.JSON(http.StatusOK, j)
	return
}

type WechatRespList struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}

// JSON 提供了系统标准JSON输出方法。
func JSONWechatList(c *gin.Context, data interface{}, totalCnt int, pageSize int) {
	totalPage := (totalCnt / pageSize) + 1

	result := new(JSONResult)
	result.Code = 0
	result.Message = "ok"
	result.Data = WechatRespList{
		List:  data,
		Total: totalPage, // 页码
	}
	c.JSON(http.StatusOK, result)
	return
}
