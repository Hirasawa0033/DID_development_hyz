package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 测试定死了code、token，真实开发应该持久化
var (
	AuthCode        = "CODE_CODE"
	AccessTokenCode = "AccessTokenCode"
)

func main() {
	r := gin.Default()

	// 前端获取 code
	r.POST("/code", Code)
	// 后端通过 code 获取令牌 token
	r.POST("/access_token", AccessToken)
	// 后端通过 token 获取用户信息
	r.POST("/user_info", UserInfo)

	r.Run(":9090")
}

type UserInfoParams struct {
	Token string `json:"token"`
}

func UserInfo(c *gin.Context) {
	params := &UserInfoParams{}
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Token != AccessTokenCode {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(200, gin.H{"user_name": "zhangsan", "Age": 20})
}

type TokenParams struct {
	Code string `json:"code"`
}

func AccessToken(c *gin.Context) {
	params := &TokenParams{}
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if params.Code != AuthCode {
		c.JSON(200, gin.H{"access_token": ""})
		return
	}

	c.JSON(200, gin.H{"access_token": AccessTokenCode})
}

type CodeParams struct {
	ClientId    string `json:"client_id"`
	Code        string `json:"code"`
	Scope       string `json:"scope"`
	State       string `json:"state"`
	RedirectUri string `json:"redirect_uri"`
}

func Code(c *gin.Context) {
	params := &CodeParams{}
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.Code = AuthCode

	c.JSON(200, gin.H{"msg": params})
}
