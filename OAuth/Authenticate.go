package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	postTokenUrl    = "http://localhost:9090/access_token"
	postUserInfoUrl = "http://localhost:9090/user_info"
)

func main() {
	r := gin.Default()

	// 获取 token
	r.POST("/token", Token)

	r.Run(":9091")
}

type TokenParams struct {
	Code string `json:"code"` // 授权码 是前端拉起第三方页面获取到的code
	// ==== 以下本次测试的用不到 ====
	GrantType   string `json:"grant_type"`   // 授权流程
	RedirectUri string `json:"redirect_uri"` // 重定向uri
}

func Token(c *gin.Context) {
	params := &TokenParams{}
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1.请求第三方颁发 token
	assToken := &AccToken{}
	Post(assToken, postTokenUrl, "code", params.Code)
	if assToken.AccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "get assecc_token error"})
		return
	}

	// 2.获取用户信息
	userInfo := &UserInfo{}
	err := Post(userInfo, postUserInfoUrl, "token", assToken.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	// ------------- 这里做一些判断操作，比如验证用户是否存在，创建/更新/登陆...... ---------------

	c.JSON(200, gin.H{
		"name": userInfo.UserName,
		"age":  userInfo.Age,
	})
}

type UserInfo struct {
	UserName string `json:"user_name"`
	Age      int    `json:"age"`
}
type AccToken struct {
	AccessToken string `json:"access_token"`
}

// 请求工具类
func Post(params interface{}, url string, key string, val string) error {
	client := &http.Client{}
	data := map[string]interface{}{
		key: val,
	}
	reqData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(reqData))
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	if len(respBody) == 0 {
		return errors.New("body nil")
	}

	json.Unmarshal(respBody, params)

	return nil
}
