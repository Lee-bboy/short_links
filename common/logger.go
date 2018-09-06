package common

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var AccessLogger *log.Logger

func init() {
	file, err := os.OpenFile("log/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}

	AccessLogger = log.New(file, "", log.Ldate|log.Ltime)
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		//计算HTTP-Referer
		referer := "-"
		if hr := r.Referer(); len(hr) != 0 {
			referer = hr
		}

		url := r.URL

		AccessLogger.Printf(
			`%s "%s %s %s" %s %s %s`,
			strings.Split(r.RemoteAddr, ":")[0],
			r.Method,
			url.RequestURI(),
			r.Proto,
			referer,
			r.UserAgent(),
			time.Since(start),
		)
	})
}
