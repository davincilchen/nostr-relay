package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//s Server
func Router(s *Server) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(Logger, gin.Recovery())

	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", SocketHandler)
	router.GET("/watcher", func(c *gin.Context) {
		c.HTML(http.StatusOK, "watcher.html", gin.H{
			"title": "Posts",
		})
	})

	router.GET("/exit", s.Shutdown)
	router.GET("/info", info)

	return router
}

func info(c *gin.Context) {
	c.JSON(200, gin.H{ // response json
		"version": "0.0.0.1",
	})
}
