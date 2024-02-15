package handlers

import (
	"github.com/ValeryVerkhoturov/chat/utils/i18nUtils"
)

type TemplateData struct {
	APIVersion  int
	Data        interface{}
	PublicUrl   string
	Locale      i18nUtils.Locale
	LocaleName  string
	TelegramUrl string
}
