package netHandler

import (
	"GoProx/config"
	"GoProx/proxError"
	"bytes"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/disintegration/gift"
	"image"
	"net/http"
	"strconv"
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
	buffer, err := GetFromCache(h.Cache, h.OriginRequest.URL.Path)
	if buffer == nil {
		//not found, fetch from remote then resize
		sourceRemotePath := h.extractSourceRemotePath()
		askedHeight, askedWidth, err := h.extractSizeFromUrl(h.OriginRequest.URL.Path)
		sourceRemotePath = fmt.Sprintf("%s%s", h.Config.RelayAdress, sourceRemotePath)
		buffer, err = GetFromCacheOrRemote(sourceRemotePath, h.Cache)
		if err != nil {
			_ = WriteBadRequest(h.Writer)
		}

		newImage, err := h.resize(buffer, askedHeight, askedWidth)
		if err != nil {
			_ = WriteBadRequest(h.Writer)
		}

		_ = h.Cache.Set(h.OriginRequest.URL.Path, *newImage)
		WriteOk(h.Writer, *newImage)

	}

	if err != nil {
		_ = WriteBadRequest(h.Writer)
	}

}

func (h *CompressHandler) extractSourceRemotePath() string {
	split := strings.Split(h.OriginRequest.URL.Path, "/")
	l := len(split)
	remotePath := strings.Join(split[2:l], "/")
	return remotePath
}

func (h *CompressHandler) resize(buffer *[]byte, width int, height int) (*[]byte, error) {
	if width < 10 || height < 10 {
		return buffer, proxError.InvalidSize //todo
	}
	//	newImg, err := bimg.NewImage(*bytes).Resize(width, height)
	//	if err != nil {
	//	return nil, err
	//}
	ioReader := bytes.NewReader(*buffer)
	srcImage, _, _ := image.Decode(ioReader)
	g := gift.New(
		gift.Resize(width, height, gift.LanczosResampling),
	)
	dst := image.NewRGBA(g.Bounds(srcImage.Bounds()))
	dstBuffer := []byte(dst.Pix)
	return &dstBuffer, nil

}

func (h *CompressHandler) extractSizeFromUrl(path string) (int, int, error) {
	split := strings.Split(h.OriginRequest.URL.Path, "/")
	strSize := split[1]
	arrHw := strings.Split(strSize, "-")
	if len(arrHw) != 2 {
		return -1, -1, proxError.InvalidSize
	}
	width, err := strconv.Atoi(arrHw[0])
	height, err := strconv.Atoi(arrHw[1])
	return width, height, err

}
