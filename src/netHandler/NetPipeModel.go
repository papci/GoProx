package netHandler

import (
	"GoProx/config"
	"github.com/allegro/bigcache/v3"
	"log"
	"net/http"
)

type NetPipeModel struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Cache   bigcache.BigCache
	Conf    config.Config
	Logger  *log.Logger
}
