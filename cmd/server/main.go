package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/justcompile/tnl/pkg/proxy"
	"github.com/justcompile/tnl/pkg/socketserver"
	log "github.com/sirupsen/logrus"
)

func main() {
	ws := gin.Default()

	hub := socketserver.NewHub()
	go hub.Run()

	ws.Any(socketserver.WebSocketPath, func(c *gin.Context) {
		socketserver.ServeWs(hub, c.Writer, c.Request)
	})

	go func(e *gin.Engine) {
		if err := endless.ListenAndServe(":8081", e); err != nil {
			log.Fatal(err)
		}
	}(ws)

	r := gin.Default()

	r.Any("/*proxyPath", gin.WrapH(&proxy.Handler{Hub: hub}))

	if err := endless.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
