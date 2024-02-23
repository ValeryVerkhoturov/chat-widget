package v1

import (
	"bytes"
	"github.com/ValeryVerkhoturov/chat/config"
	"github.com/ValeryVerkhoturov/chat/handlers"
	"github.com/ValeryVerkhoturov/chat/templates"
	"github.com/ValeryVerkhoturov/chat/utils/requestUtils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ChatWidget(c *gin.Context) {
	locale, localeName := requestUtils.GetLocale(c)

	c.Header("Content-Type", "application/javascript")

	var buf bytes.Buffer
	err := templates.HTML.ExecuteTemplate(&buf, "chat-widget.html", handlers.TemplateData{
		APIVersion:  1,
		PublicUrl:   config.PublicUrl,
		Locale:      locale,
		LocaleName:  localeName,
		TelegramUrl: config.TelegramUrl,
	})
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var jsContent = requestUtils.WrapHTMLWithEmbeddingJS(buf)

	c.String(http.StatusOK, jsContent)
}

func ChatContainer(c *gin.Context) {
	locale, localeName := requestUtils.GetLocale(c)

	c.HTML(http.StatusOK, "chat-container.html", handlers.TemplateData{
		APIVersion:  1,
		PublicUrl:   config.PublicUrl,
		Locale:      locale,
		LocaleName:  localeName,
		TelegramUrl: config.TelegramUrl,
	})
}
