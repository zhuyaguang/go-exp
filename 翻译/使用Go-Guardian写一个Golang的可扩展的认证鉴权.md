# 用Go-Guardian写一个Golang的可扩展的身份认证

> 作者： Sanad Haj  译者:朱亚光  策划：Tina
>
>  Sanad Haj：就职于F5Networks软件工程师 
>
> 原文链接 [Writing Scalable Authentication in Golang Using Go-Guardian](https://medium.com/@hajsanad/writing-scalable-authentication-in-golang-using-go-guardian-83691219a73a)

在构建web和REST API 应用中，如何打造一个用户信任和依赖的系统是非常重要的。

身份认证很重要，因为它通过只允许经过身份认证的用户访问其受保护的资源，从而使得机构和应用程序能够来保持其网络的安全。

在本教程中，我们将讨论如何使用Golang和Go-Guardian库来处理运行在集群模式下程序的身份验证。



### 问题

只要用户信息存储或者缓存在服务器端，身份认证就是一个可能会导致扩展性问题的地方。

在Kubernetes、docker swarm等集群模式下，甚至在LB后端，运行无状态的应用程序，都不能保证将单个服务器分配给特定的用户。



### 用例和解决方案

假设我没有两个可复制的应用程序A和B，并且运行在LB后面。当用户通过LB路由向应用程序A请求token，这个时候token已经产生并且缓存在应用程序中，同时同一个用户通过LB路由向应用程序B请求受保护的资源，这个会导致身份认证错误而请求失败。

让我们想想如何解决上述问题，在不降低性能的情况下扩展应用程序，并记住这种服务必须是无状态的。

**建议解决方法：**

* token存储在db中，服务器中程序缓存。
* 分布式缓存
* 共享缓存
* 粘性会话

上面所有的方法都会面临同一个问题，我们试想一下，如果数据库或者共享缓存挂了，甚至程序本身挂了会发生什么？

解决这类问题的最佳解决方案就是使用无状态token，在该token里面可以再次对其进行签名和验证。

在本教程中，我们将使用RFC 7519中定义的JWT，主要是因为其在网络上大家使用的比较广泛，都使用过是听说过。



### Go-Guardian 概述

Go-Guardian 是一个golang库，它提供了一种简单、简洁和惯用的方法来构造强大先进的API和web身份验证。



Go-Guardian的唯一目的就是验证请求，他通过一组被称为策略的可扩展的身份认证方法来实现。Go-Guardian不挂载路由也不假设任何特定的数据库模式，这极大提高了灵活性，允许开发者自己做决定。

API很简单：你提交请求给Go-Guardian进行身份验证，Go-Guardian调用策略来进行最终用户的请求认证。

策略提供回调方法来控制当身份认证成功或者失败的情况。



### 为什么要使用Go-Guardian 

当构建一个现代应用程序时，你肯定不希望重复造轮子。而且当你聚焦精力构建一个优秀的软件时，Go-Guardian正好解决了你的燃眉之急。

下面有几个可以让你尝试一下的理由：

* 提供了简单、简介、惯用的API。
* 提供了最流行和最传统的身份认证方法。
* 基于不同的机制和算法，提供一个包来缓存身份验证决策。
* 提供了基于RFC-4226和RFC-6238的双向身份认证和一次性密码。



### 创建我们的项目

我们开始新建一个项目

`mkdir scalable-auth && cd scalable-auth && go mod init scalable-auth && touch main.go`

新建了一个“scalable-auth”的文件夹，并且go.mod初始化。

当然我们也需要安装gorilla mux,、go-guardian、[jwt-go](https://github.com/dgrijalva/jwt-go)。

```go
`go get github.com/gorilla/mux`
`go get github.com/shaj13/go-guardian`
`go get "github.com/dgrijalva/jwt-go"`
```



### 我们的第一行代码

在我们写任何代码之前，我们需要写一些强制代码来运行程序。

```go
package main
import (
  "log"
)
func main() {
  log.Println("Auth !!")
}
```

### 创建我们的endpoints

我们将删掉打印“Auth!!”那行代码，添加gorilla Mux包初始化路由。

```go
package main
import (
  "github.com/gorilla/mux"
)
func main() {
  router := mux.NewRouter()
}
```

现在我们要建立我们API的endpoints,我们把所有的endpoints都创建在main函数里面，每一个endpoint都需要一个函数来处理请求。

```go
package main
import (
  "net/http"
  "log"
  "github.com/gorilla/mux"
)
func main() {
  router := mux.NewRouter() 
  router.HandleFunc("/v1/auth/token", createToken).Methods("GET")
  router.HandleFunc("/v1/book/{id}", getBookAuthor).Methods("GET")
  log.Println("server started and listening on http://127.0.0.1:8080")
  http.ListenAndServe("127.0.0.1:8080", router)
}
```

我们创建了两个路由，第一个是获取token的API,第二个是获取受保护的资源的信息，即通过id书的作者信息。

### 路由处理程序

现在我们只需要定义处理请求的函数了

**createToken()**

```go
func createToken(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,    jwt.MapClaims{
		"iss": "auth-app",
		"sub":  "medium",
		"aud": "any",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
jwtToken, _:= token.SignedString([]byte("secret"))
     w.Write([]byte(jwtToken))
}
```

**getBookAuthor()**

```go
func getBookAuthor(w http.ResponseWriter, r *http.Request) {    vars := mux.Vars(r)
    id := vars["id"]
    books := map[string]string{
        "1449311601": "Ryan Boyd",
        "148425094X": "Yvonne Wilson",
        "1484220498": "Prabath Siriwarden",
    }
    body := fmt.Sprintf("Author: %s \n", books[id])
    w.Write([]byte(body))
}
```

现在我们来发送一些简单的请求来测试下代码！

```shell
curl  -k http://127.0.0.1:8080/v1/book/1449311601 
Author: Ryan Boyd

curl  -k http://127.0.0.1:8080/v1/auth/token

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE1OTczNjE0NDYsImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.EepQzhuAS-lnljTZad3vAO2vRbgflB53aUCfCnlbku4
```

### 使用Go-Guardian集成

首先我们在main函数前面添加下面的变量定义

```go
var authenticator auth.Authenticator
var cache store.Cache
```

接着我们写两个函数来验证用户的凭证和token

```go
func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
    if userName == "medium" && password == "medium" {
        return auth.NewDefaultUser("medium", "1", nil, nil), nil
    }
    return nil, fmt.Errorf("Invalid credentials")
}
func verifyToken(ctx context.Context, r *http.Request, tokenString string) (auth.Info, error) {
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
   return nil, fmt.Errorf("Unexpected signing method: %v",     token.Header["alg"])
}
   return []byte("secret"), nil
})
if err != nil {
       return nil, err
}
if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
user := auth.NewDefaultUser(claims["medium"].(string), "", nil, nil)
return user, nil
}
return nil , fmt.Errorf("Invaled token")
}
```

我们还需要一个函数来新建Go-Guardian.

```go
func setupGoGuardian() {   
    authenticator = auth.New()
    cache = store.NewFIFO(context.Background(), time.Minute*5)
    basicStrategy := basic.New(validateUser, cache) 
    tokenStrategy := bearer.New(verifyToken, cache)
    authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
    authenticator.EnableStrategy(bearer.CachedStrategyKey,    tokenStrategy)
}
```

我们构造一个authenticator来接受请求，并且将其分发给策略，并且第一个成功验证的请求返回用户信息。另外初始化一块缓存来缓存身份认证的结果能够提高服务器性能。



接着我们需要一个HTTP的中间件来拦截请求，使得请求到达最终的路由之前进行用户的身份验证。

```go
func middleware(next http.Handler) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Executing Auth Middleware")
        user, err := authenticator.Authenticate(r)
        if err != nil {
            code := http.StatusUnauthorized
            http.Error(w, http.StatusText(code), code)
            return
        }
        log.Printf("User %s Authenticated\n", user.UserName())
        next.ServeHTTP(w, r)
    })
}
```

最后我们把createToken和getBookAuthor函数封装下，用中间件来请求身份验证。

```
middleware(http.HandlerFunc(createToken))
middleware(http.HandlerFunc(getBookAuthor))
```

不要忘记在第一个main函数之前调用下 GoGuardian 

```
setupGoGuardian()
```



### 测试下我们的API

首先我们在两个不同的shell终端里面两次运行程序

```
PORT=8080 go run main.go
PORT=9090 go run main.go
```

从副本A(8080端口)获取token

```
curl  -k http://127.0.0.1:8080/v1/auth/token -u medium:medium

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE1OTczNjI4NjksImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.SlignTJE3YD9Ecl24ygoYRu_9tVucCLop4vXWKzaRTw
```

从副本B(9090端口)使用token获取书的作者

```
curl  -k http://127.0.0.1:8080/v1/book/1449311601 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE1OTczNjI4NjksImlzcyI6ImF1dGgtYXBwIiwic3ViIjoibWVkaXVtIn0.SlignTJE3YD9Ecl24ygoYRu_9tVucCLop4vXWKzaRTw"
Author: Ryan Boyd
```



### 感谢你的阅读

希望这篇文章对你有用，至少希望能够帮助你们熟悉使用Go-Guardian来构建一个最基本的服务端身份认证。很多关于 Go-Guardian你可以访问[*GitHub* ](https://github.com/shaj13/go-guardian/tree/master/auth)and [*GoDoc*](https://pkg.go.dev/github.com/shaj13/go-guardian?tab=doc)