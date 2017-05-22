/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/ankyra/escape-registry/config"
	"github.com/ankyra/escape-registry/dao"
	"github.com/ankyra/escape-registry/handlers"
	"github.com/ankyra/escape-registry/shared"
	"github.com/ankyra/escape-registry/storage"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
)

const (
	defaultConfigFile = "/etc/escape-registry/config.json"
)

func getConfigLocation(args []string) string {
	if len(args) == 1 {
		log.Println("INFO: Using default configuration file location:", defaultConfigFile)
		return defaultConfigFile
	} else if len(args) == 2 {
		log.Println("INFO: Using configuration file location:", args[1])
		return args[1]
	}
	log.Fatalln("Error: too many arguments given. Usage: escape-registry [CONFIG_FILE]")
	return ""
}

func loadConfig(configFile string) (*config.Config, error) {
	env := os.Environ()
	if !config.PathExists(configFile) {
		log.Printf("WARN: Couldn't find configuration file '%s'. Using default configuration.", configFile)
		return config.NewConfig(env)
	} else {
		log.Printf("INFO: Loading configuration file '%s\n", configFile)
		return config.LoadConfig(configFile, env)
	}
}

func activateConfig(conf *config.Config) error {
	log.Printf("INFO: Activating '%s' database\n", conf.Database)
	if err := dao.LoadFromConfig(conf); err != nil {
		return err
	}
	log.Printf("INFO: Activating '%s' storage backend\n", conf.StorageBackend)
	if err := storage.LoadFromConfig(conf); err != nil {
		return err
	}
	return nil
}

func loadAndActivateConfig() *config.Config {
	configFile := getConfigLocation(os.Args)
	conf, err := loadConfig(configFile)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}
	if err := activateConfig(conf); err != nil {
		log.Fatalln("ERROR:", err.Error())
	}
	return conf
}

func getMux() *mux.Router {
	r := mux.NewRouter()
	getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", HomeHandler)

	getRouter.Handle("/apps/", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.RegistryHandler))))
	getRouter.Handle("/apps/{name}/", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.RegistryHandler))))
	getRouter.Handle("/apps/{name}/{version}/", negroni.New(
		negroni.Wrap(http.HandlerFunc(handlers.RegistryHandler))))

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
	return r
}

func getHandler(router *mux.Router) http.Handler {
	middleware := negroni.New()
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	middleware.Use(recovery)
	middleware.Use(negroni.NewLogger())
	middleware.UseHandler(router)
	return middleware
}

func StartRegistry(router *mux.Router) {
	fmt.Println(shared.EscapeLogo)
	config := loadAndActivateConfig()

	handler := getHandler(router)
	http.Handle("/", handler)

	port := config.Port
	log.Printf("INFO: Starting Escape Registry v%s on port %s\n", registryVersion, port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func main() {
	StartRegistry(getMux())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Escape Release Registry v" + registryVersion))
}
