package main

import (
	"net/http"
	"server/syslog"
	"server/util"
	"strconv"
)

func main() {
	router()
}

func router() {
	http.HandleFunc("/", Serve)
	http.HandleFunc("/static/", ServeStatic)
	syslog.Info("main   Server is started from :" + strconv.Itoa(util.Config.Port))
	syslog.Fatal(http.ListenAndServe(util.Config.Ip+":"+strconv.Itoa(util.Config.Port), nil))
}
