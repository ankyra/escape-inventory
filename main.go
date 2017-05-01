package main

import (
    "os"
	"log"
	"net/http"
	"github.com/ankyra/escape-registry/config"
	"github.com/ankyra/escape-registry/handlers"
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/storage"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	defaultConfigFile = "config.json"
)

func loadConfig(configFile string) (*config.Config, error) {
    env := os.Environ()
    if !config.PathExists(configFile) {
        log.Println("Using default configuration")
        return config.NewConfig(env)
    } else {
        log.Printf("Loading configuration file '%s\n", configFile)
        return config.LoadConfig(configFile, env)
    }
}

func activateConfig(conf *config.Config) error {
    log.Printf("Activating '%s' database\n", conf.Database)
    if err := dao.LoadFromConfig(conf); err != nil {
        return err
    }
    log.Printf("Activating '%s' storage backend\n", conf.StorageBackend)
    if err := storage.LoadFromConfig(conf); err != nil {
        return err
    }
    return nil
}

func main() {
    conf, err := loadConfig(defaultConfigFile)
	if err != nil {
        log.Fatalln("Error:", err.Error())
	}
    if err := activateConfig(conf); err != nil {
        log.Fatalln("Error:", err.Error())
	}

	r := mux.NewRouter()
    getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", HomeHandler)
	getRouter.Handle("/r/{release}/", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.GetMetadataHandler))))
	getRouter.Handle("/r/{release}/download", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.DownloadHandler))))
	getRouter.Handle("/r/{release}/next-version", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.NextVersionHandler))))
	getRouter.Handle("/export-releases", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.ExportReleasesHandler))))

    postRouter := r.Methods("POST").Subrouter()
	postRouter.Handle("/r/", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.RegisterHandler)))).Methods("POST")
	postRouter.Handle("/r/{release}/upload", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.UploadHandler))))
	postRouter.Handle("/import-releases", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.ImportReleasesHandler))))

	middleware := negroni.New()
    recovery := negroni.NewRecovery()
    recovery.PrintStack = false
    middleware.Use(recovery)
    middleware.Use(negroni.NewLogger())
	middleware.UseHandler(r)
	http.Handle("/", middleware)

	port := "3000"
	log.Printf("Starting Escape Registry v%s on port %s\n", registryVersion, port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Escape Release Registry v" + registryVersion))
}
