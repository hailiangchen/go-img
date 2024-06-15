package server

import (
	"github.com/gin-gonic/gin"
	"go-img/conf"
	"go-img/controller"
	"net/http"
)

func Run() {
	r := gin.Default()
	//r.SetTrustedProxies(nil)

	gin.SetMode(conf.AppConf.HttpConf.RunMode)

	r.StaticFile("/favicon.ico", "./www/favicon.ico")

	admin := new(controller.AdminController)
	admin.RegisterRouter(r, authFunc())

	index := new(controller.IndexController)
	index.RegisterRouter(r)

	r.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "404")
	})

	r.Static("/admin/index/", "./www")

	//限制上传最大尺寸
	r.MaxMultipartMemory = 8 << 20
	r.Run(conf.AppConf.HttpConf.Address)
}
