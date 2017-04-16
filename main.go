package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	defaultConfigFile = "config.json"
)

var (
	cfg      *Config
	store    *sessions.CookieStore
)

func main() {
	var err error
	cfg, err = loadConfig(defaultConfigFile)
	if err != nil {
		panic(err)
	}

	store = sessions.NewCookieStore([]byte(cfg.Secret))

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)

    port := "3000"
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("test"))
}
