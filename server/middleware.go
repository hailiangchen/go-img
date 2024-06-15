package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-img/util"
	"net"
	"net/http"
	"strings"
)

func authFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteIP := getRemoteIp(c.Request)
		if util.IsAllow(remoteIP) {
			c.Next()
		} else {
			c.AbortWithError(http.StatusForbidden, errors.New("拒绝访问"))
		}
	}
}

func getRemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = strings.Split(ip, ",")[0]
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
