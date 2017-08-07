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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var ReadRoutes = map[string]http.HandlerFunc{
	"/":       HomeHandler,
	"/health": handlers.HealthCheckHandler,

	"/api/v1/registry/":                                                           handlers.RegistryHandler,
	"/api/v1/registry/{project}/":                                                 handlers.RegistryHandler,
	"/api/v1/registry/{project}/units/":                                           handlers.RegistryHandler,
	"/api/v1/registry/{project}/units/{name}/":                                    handlers.RegistryHandler,
	"/api/v1/registry/{project}/units/{name}/versions/":                           handlers.RegistryHandler,
	"/api/v1/registry/{project}/units/{name}/versions/{version}/":                 handlers.RegistryHandler,
	"/api/v1/registry/{project}/units/{name}/versions/{version}/diff/":            handlers.DiffHandler,
	"/api/v1/registry/{project}/units/{name}/versions/{version}/diff/{diffWith}/": handlers.DiffHandler,
	"/api/v1/registry/{project}/units/{name}/versions/{version}/download":         handlers.DownloadHandler,
	"/api/v1/registry/{project}/units/{name}/versions/{version}/previous/":        handlers.PreviousHandler,
	"/api/v1/registry/{project}/units/{name}/next-version":                        handlers.NextVersionHandler,

	"/api/v1/internal/export": handlers.ExportReleasesHandler,
}

var WriteRoutes = map[string]http.HandlerFunc{
	"/api/v1/registry/{project}/upload":                                 handlers.RegisterAndUploadHandler,
	"/api/v1/registry/{project}/register":                               handlers.RegisterHandler,
	"/api/v1/registry/{project}/units/{name}/versions/{version}/upload": handlers.UploadHandler,

	"/api/v1/internal/import": handlers.ImportReleasesHandler,
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Escape Release Registry [+project] v" + cmd.RegistryVersion))
}

func getMux() *mux.Router {
	r := mux.NewRouter()
	getRouter := r.Methods("GET").Subrouter()
	for url, handler := range ReadRoutes {
		getRouter.Handle(url, handler)
	}
	postRouter := r.Methods("POST").Subrouter()
	for url, handler := range WriteRoutes {
		postRouter.Handle(url, handler)
	}
	r.Handle("/metrics", promhttp.Handler())
	return r
}

func main() {
	cmd.StartRegistry(getMux())
}
