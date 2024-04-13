package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
)

func main() {
	// 1. 配置OAuth 2.0客户端
	conf := &oauth2.Config{
		ClientID:     "1145141414510",
		ClientSecret: "18428715198419594917448747178174",
		Scopes:       []string{"scope1", "scope2"},
		RedirectURL:  "http://localhost:8000/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://authserver.com/auth",
			TokenURL: "https://authserver.com/token",
		},
	}

	// 2. 获取授权码
	authURL := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	var code string
	fmt.Print("Enter the authorization code: ")
	fmt.Scan(&code)

	// 3. 通过授权码获取访问令牌
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Failed to exchange token: %v", err)
	}

	// 4. 使用访问令牌访问受保护资源
	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://protected-resource.com/api")
	if err != nil {
		log.Fatalf("Failed to get resource: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	// ...
}
