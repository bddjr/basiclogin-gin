package basiclogin

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ScriptRedirect(ctx *gin.Context, code int, path string) {
	ctx.Data(code, "text/html; charset=utf-8", []byte(`<script>location.replace(`+strconv.Quote(path)+`)</script>`))
}

// ⚠ If you need *http.Cookie, please use
//
//	ctx.Writer.Header().Add("Set-Cookie", cookie.String())
func New(group *gin.RouterGroup, callBack func(ctx *gin.Context, username string, password string, secure bool)) {
	const timeBase = 36
	const cookieName = "BasicLoginTime"

	redirect := func(ctx *gin.Context) {
		ScriptRedirect(ctx, 401, group.BasePath()+"/"+strconv.FormatInt(time.Now().UnixMilli(), timeBase))
	}

	group.Use(func(ctx *gin.Context) {
		ctx.Header("Referrer-Policy", "same-origin")
	})

	group.GET("/", redirect)

	group.GET("/:t", func(ctx *gin.Context) {
		param := ctx.Param("t")

		// 如果上次登录的时候用过这个时间戳，重新生成网址。
		// If this timestamp was used during the last login, regenerate the URL.
		// Fix for firefox.
		cookieBasicLoginUsed, _ := ctx.Cookie(cookieName)
		if cookieBasicLoginUsed == param {
			redirect(ctx)
			return
		}

		paramTimeInt, err := strconv.ParseInt(param, timeBase, 64)
		if err != nil {
			redirect(ctx)
			return
		}
		// 如果cookie记录的时间戳比地址栏记录的新，重新生成网址。
		// If the timestamp recorded in the cookie is newer than the one recorded in the address bar, regenerate the URL.
		// Fix for firefox.
		if cookieBasicLoginUsed != "" {
			cookieTimeInt, err := strconv.ParseInt(cookieBasicLoginUsed, timeBase, 64)
			if err == nil && cookieTimeInt > paramTimeInt {
				redirect(ctx)
				return
			}
		}
		// 时间戳不能超过当前服务端的时间戳。
		// The timestamp cannot exceed the timestamp of the current server
		paramTime := time.UnixMilli(paramTimeInt)
		if paramTime.After(time.Now()) {
			redirect(ctx)
			return
		}

		referer := ctx.Request.Header.Get("Referer")
		if referer == "" {
			if paramTime.Add(3 * time.Second).After(time.Now()) {
				// 如果网址时间戳对比当前时间戳，相差不超过3秒，那么浏览器不支持Referer。
				// If the timestamp of the website differs from the current timestamp by no more than 3 seconds, then the browser does not support Referer.
				ctx.String(400, "Missing Referer Header")
				return
			}
			redirect(ctx)
			return
		}
		secure := strings.HasPrefix(referer, "https")

		username, password, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.Header("WWW-Authenticate", `Basic realm=`+ctx.Request.URL.Path+`, charset="UTF-8"`)
			ctx.Status(401)
			return
		}
		ctx.SetCookie(cookieName, param, 0, group.BasePath(), "", secure, true)
		callBack(ctx, username, password, secure)
	})
}
