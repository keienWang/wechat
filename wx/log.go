package wx

import (
	"fmt"
	"github.com/keienWang/wechat/config"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func WriteLog(r *http.Request, t time.Time, match string, pattern string) {

	if config.LogLevel != "prod" {

		d := time.Now().Sub(t)

		l := fmt.Sprintf("[ACCESS] | % -10s | % -40s | % -16s | % -10s | % -40s |", r.Method, r.URL.Path, d.String(), match, pattern)

		log.Println(l)
	}
}
