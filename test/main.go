package main

import (
	"crypto/hmac"
	"fmt"

	"github.com/bddjr/basiclogin-gin"
	"github.com/gin-gonic/gin"
)

const staticUserName = "test"
const StaticPassword = "123456"

const cookieName = "test"
const cookieValue = "123456"

func main() {
	Router := gin.New()

	// web
	Router.GET("/", func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(cookieName)
		if err != nil || !hmac.Equal([]byte(cookie), []byte(cookieValue)) {
			basiclogin.ScriptRedirect(ctx, 400, "/login/")
			return
		}
		ctx.File("index.html")
	})

	// login
	loginGroup := Router.Group("login")

	loginGroup.Use(func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(cookieName)
		if err == nil && hmac.Equal([]byte(cookie), []byte(cookieValue)) {
			basiclogin.ScriptRedirect(ctx, 400, "/")
			ctx.Abort()
		}
	})

	basiclogin.New(loginGroup, func(ctx *gin.Context, username, password string, secure bool) {
		// ⚠ If you need *http.Cookie, please use
		//   ctx.Writer.Header().Add("Set-Cookie", cookie.String())
		if username == staticUserName && hmac.Equal([]byte(password), []byte(StaticPassword)) {
			ctx.SetCookie(cookieName, cookieValue, 0, "/", "", secure, true)
			basiclogin.ScriptRedirect(ctx, 200, "/")
			return
		}
		ctx.Status(401)
	})

	// logout
	Router.GET("/logout", func(ctx *gin.Context) {
		ctx.SetCookie(cookieName, "x", -1, "", "", false, true)
		ctx.Redirect(303, "/login/")
	})

	// listen
	fmt.Print("\n  http://localhost:8080\n\n")
	err := Router.Run(":8080")
	fmt.Println(err)
}
