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
	"net/http"
)

var ReadRoutes = map[string]http.HandlerFunc{
	"/":                         HomeHandler,
	"/apps/":                    handlers.RegistryHandler,
	"/apps/{name}/":             handlers.RegistryHandler,
	"/apps/{name}/{version}/":   handlers.RegistryHandler,
	"/r/{release}/":             handlers.GetMetadataHandler,
	"/r/{release}/download":     handlers.DownloadHandler,
	"/r/{release}/next-version": handlers.NextVersionHandler,
	"/export-releases":          handlers.ExportReleasesHandler,
}

var WriteRoutes = map[string]http.HandlerFunc{
	"/r/": handlers.RegisterHandler,
	"/r/{release}/upload": handlers.UploadHandler,
	"/import-releases":    handlers.ImportReleasesHandler,
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Escape Release Registry v" + cmd.RegistryVersion))
}

func getMux() *mux.Router {
	r := mux.NewRouter()
	getRouter := r.Methods("GET").Subrouter()
	for url, handler := range ReadRoutes {
		getRouter.HandleFunc(url, handler)
	}
	postRouter := r.Methods("POST").Subrouter()
	for url, handler := range WriteRoutes {
		postRouter.HandleFunc(url, handler)
	}
	return r
}

func main() {
	cmd.StartRegistry(getMux())
}
