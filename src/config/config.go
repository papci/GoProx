package config

import (
	"GoProx/proxError"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port        int32  `json:"port"`
	ListenIface string `json:"interface"`
	RelayAdress string `json:"relay"`
}

func FromFile(s string) Config {
	strData, err := os.ReadFile(s)
	proxError.PanicIfErr(err, log.Default())

	bytes := []byte(strData)

	var c Config
	err = json.Unmarshal(bytes, &c)
	proxError.PanicIfErr(err, log.Default())
	return c
}
