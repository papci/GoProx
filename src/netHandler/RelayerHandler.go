package netHandler

import (
	"GoProx/config"
	"github.com/allegro/bigcache/v3"
	"net/http"
)

type RelayerHandler struct {
	Origin        string
	Config        *config.Config
	Writer        http.ResponseWriter
	OriginRequest *http.Request
	Cache         *bigcache.BigCache
}

func (handler *RelayerHandler) HandleRequest() {
	bytes, err := GetFromCacheOrRemote(handler.Config.RelayAdress+handler.OriginRequest.URL.Path, handler.Cache)
	if err != nil {
		_ = WriteBadRequest(handler.Writer)
	}

	handler.Writer.WriteHeader(200)
	_, _ = handler.Writer.Write(*bytes)
}
