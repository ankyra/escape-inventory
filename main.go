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

func getMux() *mux.Router {
	r := mux.NewRouter()
	getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", HomeHandler)

	getRouter.HandleFunc("/apps/", handlers.RegistryHandler)
	getRouter.HandleFunc("/apps/{name}/", handlers.RegistryHandler)
	getRouter.HandleFunc("/apps/{name}/{version}/", handlers.RegistryHandler)

	getRouter.HandleFunc("/r/{release}/", handlers.GetMetadataHandler)
	getRouter.HandleFunc("/r/{release}/download", handlers.DownloadHandler)
	getRouter.HandleFunc("/r/{release}/next-version", handlers.NextVersionHandler)
	getRouter.HandleFunc("/export-releases", handlers.ExportReleasesHandler)

	postRouter := r.Methods("POST").Subrouter()
	postRouter.HandleFunc("/r/", handlers.RegisterHandler)
	postRouter.HandleFunc("/r/{release}/upload", handlers.UploadHandler)
	postRouter.HandleFunc("/import-releases", handlers.ImportReleasesHandler)
	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Escape Release Registry v" + cmd.RegistryVersion))
}

func main() {
	cmd.StartRegistry(getMux())
}
