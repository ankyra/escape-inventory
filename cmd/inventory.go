/*
Copyright 2017, 2018 Ankyra

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

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ankyra/escape-inventory/config"
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/metrics"
	"github.com/ankyra/escape-inventory/model"
	"github.com/ankyra/escape-inventory/storage"
	basicauth "github.com/aphistic/negroni-basicauth"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var Config *config.Config

const (
	defaultConfigFile = "/etc/escape-inventory/config.json"
)

func getConfigLocation(args []string) string {
	if len(args) == 1 {
		log.Println("INFO: Using default configuration file location:", defaultConfigFile)
		return defaultConfigFile
	} else if len(args) == 2 {
		log.Println("INFO: Using configuration file location:", args[1])
		return args[1]
	}
	log.Fatalln("Error: too many arguments given. Usage: escape-inventory [CONFIG_FILE]")
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
	log.Printf("INFO: Updating unprocessed release dependencies\n")
	if err := model.ProcessUnprocessedReleases(); err != nil {
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

type MetricMiddleware struct{}

func NewMetricMiddleware() *MetricMiddleware {
	return &MetricMiddleware{}
}

func (m *MetricMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	elapsed := time.Since(start)

	status := strconv.Itoa(res.Status())
	metrics.ResponsesTotal.WithLabelValues(status, r.Method).Inc()
	metrics.ResponsesLatency.WithLabelValues(status, r.Method).Observe(float64(elapsed.Seconds()))
}

func GetHandler(router *mux.Router) http.Handler {
	middleware := negroni.New()
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	middleware.Use(NewMetricMiddleware())
	middleware.Use(recovery)
	middleware.Use(negroni.NewLogger())
	if Config != nil && Config.BasicAuthPassword != "" {
		log.Printf("INFO: Enabling basic authentication.\n")
		users := map[string]string{}
		users[Config.BasicAuthUsername] = Config.BasicAuthPassword
		middleware.Use(basicauth.BasicAuth("escape", users))
	}
	middleware.UseHandler(router)

	return middleware
}

func LoadConfig() *config.Config {
	fmt.Println(EscapeLogo)
	for _, arg := range os.Args[1:] {
		if arg == "--version" {
			os.Exit(0)
		}
	}
	Config = loadAndActivateConfig()
	return Config
}

func StartInventory(router *mux.Router) {
	handler := GetHandler(router)
	http.Handle("/", handler)

	port := Config.Port
	log.Printf("INFO: Starting Escape Inventory v%s on port %s\n", InventoryVersion, port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
