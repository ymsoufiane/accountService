package middlewares

import (
	"net/http"
	"account/context"
)

type http_function func(w http.ResponseWriter, r *http.Request)

type Logs struct {
	next http_function
}

func (l *Logs) LogsMiddelware(next http.Handler) http.Handler {
	l.next = next.ServeHTTP
	return l
}

func (l *Logs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.InfoLogger.Println("Ip:" + r.RemoteAddr + " Path: " + r.URL.Path + " Method:" + r.Method)
	l.next(w, r)

}

