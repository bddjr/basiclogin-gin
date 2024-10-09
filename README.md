# Warning
不建议使用，它不能记住密码。而且它会导致Firefox在不关闭浏览器的情况下无法安全退出登录。  

It is not recommended to use it, as it does not remember passwords and it can cause Firefox to be unable to securely log out without closing the browser.

# basiclogin-gin

为Cookie设计的轻量级登录框架，兼容chrome 、firefox浏览器，暂未测试safari浏览器。  

Lightweight login framework designed for cookie, compatible with chrome, firefox, not tested safari.  

Reference https://developer.mozilla.org/docs/Web/HTTP/Authentication

## Get
```
go get github.com/bddjr/basiclogin-gin
```

## Example
[See test/main.go](test/main.go)  

```go
loginGroup := Router.Group("login")
basiclogin.New(loginGroup, func(ctx *gin.Context, username, password string, secure bool) {
    // ⚠ If you need *http.Cookie, please use
    //   ctx.Writer.Header().Add("Set-Cookie", cookie.String())
    if username == staticUserName && hmac.Equal([]byte(password), []byte(StaticPassword)) {
        ctx.SetCookie(cookieName, cookieValue, 0, "/", "", secure, true)
        ctx.Header("Referrer-Policy", "no-referrer")
        basiclogin.ScriptRedirect(ctx, 401, "/")
        return
    }
	ctx.String(401, "Wrong usename or password")
})
```

## License
[BSD-3-clause license](LICENSE.txt)  
