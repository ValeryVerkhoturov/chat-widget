package handlers

import (
	"github.com/ValeryVerkhoturov/chat/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusNotFound, "index.html", TemplateData{
		APIVersion: 1,
		PublicUrl:  config.PublicUrl,
	})
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}
