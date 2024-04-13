package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()

	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	}) //这里是 /token 路由的实现

	//这非常简单。它将请求和响应传递给适当的处理程序，以便服务器可以解码请求中的所有必要的数据

	http.HandleFunc("/credentials", func(w http.ResponseWriter, r *http.Request) {
		clientId := uuid.New().String()[:8]
		clientSecret := uuid.New().String()[:8]
		err := clientStore.Set(clientId, &models.Client{
			ID:     clientId,
			Secret: clientSecret,
			Domain: "http://localhost:9094",
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"CLIENT_ID": clientId, "CLIENT_SECRET": clientSecret})
	}) //credentials 用于颁发客户端凭据 （客户端 ID 和客户端密钥）

	//它创建了两个随机字符串，一个就是客户端 ID，另一个就是客户端密钥。并把它们保存到客户端存储。然后就会返回响应
	//在这里我们使用了内存存储，但我们同样可以把它们存储到 redis，mongodb，postgres 等等里面

	http.HandleFunc("/protected", validateToken(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm protected"))
	}, srv)) //这里我们创建了一个管理器，用于客户端存储和鉴权服务本身

	log.Fatal(http.ListenAndServe(":9096", nil)) //运行这个服务并且发送 Get 请求到 http://localhost:9096/protected
}

func validateToken(f http.HandlerFunc, srv *server.Server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f.ServeHTTP(w, r)
	})
}

//运行这个代码并到 http://localhost:9096/credentials 路由去注册并获取客户端 ID 和客户端密钥。
//现在去到这个链接 http://localhost:9096/token?grant_type=client_credentials&client_id=c843ef0b&client_secret=766fe701&scope=all
//你可以得到具有过期时间和一些其他信息的授权令牌。
//现在我们得到了我们的授权令牌。但是我们的 /protected 路由依然没有被保护。我们需要设置一个方法来检查每个客户端的请求是否都带有有效的令牌。如果是的，我们就可以给予这个客户端授权。反之就不能给予授权。
//现在运行服务并在 URL 不带有 访问令牌 的情况下访问 /protected 接口。或者尝试使用错误的 访问令牌。在这两种方式下鉴权服务都会阻止你。
//现在再次从服务器获得认证信息 and 访问令牌 并发送请求到受保护的接口：
//http://localhost:9096/test?access_token=YOUR_ACCESS_TOKEN
