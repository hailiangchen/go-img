package route

import (
	"go-image/controller"
	"go-image/server"
)

// RegisterRoute Register Route.
func RegisterRoute() {
	server.HandleFunc("/", controller.Index)
	server.AuthMiddlewareHandler("/admin/upload", controller.Upload)
	server.AuthMiddlewareHandler("/admin/delete/", controller.Delete)
	server.AuthMiddlewareHandler("/admin/getall", controller.GetAll)
}
