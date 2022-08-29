package netHandler

import (
	"GoProx/config"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"net/http"
	"strings"
)

type CompressHandler struct {
	Width         int32
	Height        int32
	Origin        string
	Config        *config.Config
	Writer        http.ResponseWriter
	OriginRequest *http.Request
	Cache         *bigcache.BigCache
}

func (h *CompressHandler) HandleRequest() {
	//try to find already resized image in cache
	bytes, err := GetFromCache(h.Cache, h.OriginRequest.URL.Path)
	if bytes == nil {
		//not found fetch from remote then resize
		remotePath := h.extractRemotePath()
		remotePath = fmt.Sprintf("%s%s", h.Config.RelayAdress, remotePath)
		bytes, err = GetFromCacheOrRemote(remotePath, h.Cache)
		if err != nil {
			_ = WriteBadRequest(h.Writer)
		}

	}

	if err != nil {
		_ = WriteBadRequest(h.Writer)
	}

}

func (h *CompressHandler) extractRemotePath() string {
	split := strings.Split(h.OriginRequest.URL.Path, "/")
	l := len(split)
	remotePath := strings.Join(split[2:l], "/")
	return remotePath
}

func (h *CompressHandler) resize(bytes []byte, height int32, width int32) *[]byte {
	return &bytes //todo
}
