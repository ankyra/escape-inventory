package main

import (
	"github.com/ankyra/escape-registry/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

const (
	defaultConfigFile = "config.json"
)

var (
	cfg   *Config
)

func main() {
	var err error
	cfg, err = loadConfig(defaultConfigFile)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.Handle("/r/", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.RegisterHandler))))
	r.Handle("/r/{release}/", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.GetMetadataHandler))))
	r.Handle("/r/{release}/download", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.DownloadHandler))))
	r.Handle("/r/{release}/upload", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.UploadHandler))))
	r.Handle("/r/{release}/next-version", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.NextVersionHandler))))
	r.Handle("/export-releases", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.ExportReleasesHandler))))
	r.Handle("/import-releases", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.ImportReleasesHandler))))

	middleware := negroni.Classic()
	middleware.UseHandler(r)
	http.Handle("/", middleware)

	port := "3000"
	log.Printf("Starting Escape Registry v%s on port %s\n", registryVersion, port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Escape Release Registry v" + registryVersion))
}
