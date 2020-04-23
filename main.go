package main

import (
	"fmt"
	"github.com/leeeboo/wechat/config"
	"github.com/leeeboo/wechat/wx"
	"log"
	"net/http"
	"time"
)



func main() {
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", config.Port),
		Handler:        &wx.HttpHandler{},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	log.Println(fmt.Sprintf("Listen: %d", config.Port))
	log.Fatal(server.ListenAndServe())
}

