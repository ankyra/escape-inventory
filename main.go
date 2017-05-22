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
	"github.com/ankyra/escape-registry/cmd"
	"github.com/ankyra/escape-registry/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Escape Release Registry v" + cmd.RegistryVersion))
}

func main() {
	cmd.StartRegistry(getMux())
}
