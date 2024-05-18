package server

import (
	"go-image/config"
	"go-image/utils"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

var serveMux *http.ServeMux = http.NewServeMux()

// HandleFunc Register route from HandleFunc.
func HandleFunc(pattern string, handler http.HandlerFunc) {
	serveMux.HandleFunc(pattern, handler)
}

// Handle Register route from Handle.
func Handle(pattern string, handler http.Handler) {
	serveMux.Handle(pattern, handler)
}

func AuthMiddlewareHandler(pattern string, handler http.HandlerFunc) {
	serveMux.Handle(pattern, middlewareHandler(authFunc, handler))
}

func authFunc(w http.ResponseWriter, r *http.Request) {
	remoteIP := getRemoteIp(r)
	if !utils.IsAllow(remoteIP) {
		http.Error(w, "禁止访问", http.StatusForbidden)
		return
	}
}

func middlewareHandler(middleware http.HandlerFunc, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware(w, r)

		next.ServeHTTP(w, r)
	})
}

func getRemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// RunServer start HTTP server.
func RunServer() {
	webPath := config.GetSetting("http.webPath")

	if len(webPath) != 0 {
		serveMux.Handle("/admin/index/", http.StripPrefix("/admin/index/", http.FileServer(http.Dir(webPath))))
	}

	readTimeout, err := strconv.Atoi(config.GetSetting("http.readTimeout"))
	if err != nil {
		readTimeout = 0
	}

	writeTimeout, err := strconv.Atoi(config.GetSetting("http.writeTimeout"))
	if err != nil {
		writeTimeout = 0
	}

	serv := &http.Server{
		Addr:         config.GetSetting("http.addr"),
		Handler:      serveMux,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,  // 读超时
		WriteTimeout: time.Duration(writeTimeout) * time.Second, // 写超时
		//ReadHeaderTimeout: 5 * time.Second,
		//IdleTimeout:       5 * time.Second,
	}

	serv.SetKeepAlivesEnabled(true)
	err = serv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
