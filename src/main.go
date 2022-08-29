package main

import (
	config2 "GoProx/config"
	"GoProx/netHandler"
	"GoProx/proxError"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Default().Println("Welcome to GoProx, init cache ...")
	bc, _ := initCache()
	initWebServer(bc)

}

func initCache() (*bigcache.BigCache, error) {
	return bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

}

func initWebServer(bc *bigcache.BigCache) {
	config := config2.FromFile("config.json")
	log.Default().Println(fmt.Sprintf("Starting main webserver on %s:%d", config.ListenIface, config.Port))

	http.HandleFunc(
		"/compress",
		func(w http.ResponseWriter, r *http.Request) {
			compressHandler := &netHandler.CompressHandler{
				Writer:        w,
				Config:        &config,
				OriginRequest: r,
				Cache:         bc,
				Width:         0,
				Height:        0,
			}

			compressHandler.HandleRequest()
		},
	)

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			relayerHandler := &netHandler.RelayerHandler{
				Writer:        w,
				Config:        &config,
				OriginRequest: r,
				Cache:         bc,
			}

			relayerHandler.HandleRequest()
		})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenIface, config.Port), nil)
	proxError.PanicIfErr(err, log.Default())

}
