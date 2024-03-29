// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"short_url_go/controllers"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web/context"
	cors "github.com/beego/beego/v2/server/web/filter/cors"

	"github.com/beego/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
)

//https://www.cnblogs.com/zdz8207/p/golang-learn-7.html swagger样式以及auth swagger

// 过滤器all
var FilterToken = func(ctx *context.Context) {
	// ctx.Output.Header("Access-Control-Allow-Origin", "*")
	// logs.Info("current router path is ", ctx.Request.RequestURI)
	if ctx.Request.RequestURI != "/api/users/login" &&
		ctx.Request.RequestURI != "/api/users/register" &&
		ctx.Request.RequestURI != "/api/users/tocken/account" &&
		ctx.Request.RequestURI[0:4] == "/api" {
		//没有token
		if ctx.Input.Header("authorization") == "" {
			logs.Error("without token, unauthorized !!")
			ctx.ResponseWriter.WriteHeader(401)
			ctx.ResponseWriter.Write([]byte("no permission")) //没有权限
			return
		} else {
			//accessToken错误
			token := ctx.Input.Header("authorization")
			token = strings.Split(token, " ")[1]
			logs.Info(" accessToken错误 curernttoken: ", token)
			ok := controllers.AuthenticationJWT(token)
			if !ok {
				ctx.ResponseWriter.WriteHeader(401)
				ctx.ResponseWriter.Write([]byte("no permission"))
				return
			}
		}
	}
}

var FilterHeader = func(ctx *context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", strings.Join(allDomain, ","))
}

var allDomain []string

func init() {
	iniConf, err := config.NewConfig("ini", "../data/go/conf/secret.conf")
	if err != nil {
		panic(err)
	}
	corsDomains, err := iniConf.String("FilterHttp::AllowOrigins")
	if err != nil {
		panic(err)
	}
	allDomain = strings.Split(corsDomains, ",")
	// beego.InsertFilter("/*", beego.BeforeRouter, cors.Allow(&cors.Options{
	// 	// 允许访问所有源
	// 	AllowAllOrigins: true,
	// 	// 可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
	// 	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	// 指的是允许的Header的种类
	// 	AllowHeaders: []string{"Origin", "authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	// 	// 公开的HTTP标头列表
	// 	ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	// 	// 如果设置，则允许共享身份验证凭据，例如cookie

	// }))
	beego.InsertFilter("*", beego.BeforeRouter, FilterHeader)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowOrigins:  allDomain,
		AllowHeaders:  []string{"Origin", "authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	}))
	beego.InsertFilter("*", beego.BeforeRouter, FilterToken)
	// beego.Router("/t/:shortURL", &controllers.RedirectController{})
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/users",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/shorts",
			beego.NSInclude(
				&controllers.ShortController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
