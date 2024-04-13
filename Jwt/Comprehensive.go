package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var Auth string = ""

// 定义过期时间
const TokenExpireDuration = time.Hour * 2

// 定义secret
var MySecret = []byte("这是一段用于生成token的密钥")

// 生成jwt
func GenToken(username string) (string, error) {
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "my-project",
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString(MySecret)
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 基于JWT认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := Auth
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中的auth为空",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中的auth格式错误",
			})
			//阻止调用后续的函数
			c.Abort()
			return
		}
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的token",
			})
			c.Abort()
			return
		}
		//将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		//后续的处理函数可以通过c.Get("username")来获取请求的用户信息
		c.Next()
	}

}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authHandler(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}

	user.Username = "cyl"
	user.Password = "123456"

	if user.Username == "cyl" && user.Password == "123456" {
		//生成token
		tokenString, _ := GenToken(user.Username)
		Auth = "Bearer " + tokenString
		c.JSON(http.StatusOK, gin.H{
			"code":          200,
			"msg":           "success",
			"data":          gin.H{"token": tokenString},
			"Authorization": Auth,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}

func homeHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"username": username},
	})
}

func main() {
	r := gin.Default()
	r.GET("/auth", authHandler)
	//r.POST("/auth", authHandler) // 改为 POST 请求
	r.GET("/home", JWTAuthMiddleware(), homeHandler) // 改为 GET 请求
	r.Run(":9000")
}

//http://localhost:9000/auth
//http://localhost:9000/home
/*
总之，GET 请求用于获取资源，而 POST 请求用于提交数据
*/
//最后还是魔改了一番。乐
