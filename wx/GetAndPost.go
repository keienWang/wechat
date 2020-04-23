package wx

import (
	"github.com/leeeboo/wechat/config"
	"log"
	"net/http"
)



func Get(w http.ResponseWriter, r *http.Request) {

	client, err := NewClient(r, w, config.Token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	if len(client.Query.Echostr) > 0 {
		w.Write([]byte(client.Query.Echostr))
		return
	}

	w.WriteHeader(403)
	return
}

func Post(w http.ResponseWriter, r *http.Request) {

	client, err := NewClient(r, w, config.Token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	client.Run()
	return
}


