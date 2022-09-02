package netHandler

import (
	"GoProx/config"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"io"
	"net/http"
	"strings"
)

type IHandler interface {
	HandleRequest()
}

func RelocateToSourceUrl(path string, c *config.Config) string {
	split := strings.Split(path, "/")
	l := len(split)
	return c.RelayAdress + strings.Join(split[2:l], "/")
}

func GetFromCacheOrRemote(fullRemoteUrl string, cache *bigcache.BigCache) (*[]byte, error) {
	bytes, err := GetFromCache(cache, fullRemoteUrl)
	if err != nil {
		return nil, err
	}

	if bytes == nil || len(*bytes) < 0 {
		// no data in cache, get it from remote
		bytes, err = GetFromRemote(fullRemoteUrl)
		if err != nil {
			return nil, err
		}

		//store in cache
		err := cache.Set(fullRemoteUrl, *bytes)
		if err != nil {
			return nil, err
		}

	}

	return bytes, nil
}

func GetFromCache(cache *bigcache.BigCache, key string) (*[]byte, error) {
	bytes, err := cache.Get(key)

	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &bytes, nil
}

func GetFromRemote(fullRemoteUrl string) (*[]byte, error) {
	resp, err := http.Get(fullRemoteUrl)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &bytes, nil

}

func WriteBadRequest(writer http.ResponseWriter) error {
	writer.WriteHeader(400)
	_, err := writer.Write([]byte("Bad Request"))
	if err != nil {
		return err
	}

	return nil
}

func WriteOk(writer http.ResponseWriter, bytes []byte) {
	writer.WriteHeader(200)
	_, _ = writer.Write(bytes)
}

func ComposeRemoteUrl(config *config.Config, path string) string {
	return fmt.Sprintf("%s%s", config.RelayAdress, path)
}
