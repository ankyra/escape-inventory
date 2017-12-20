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

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ankyra/escape-inventory/cmd"
	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
	"github.com/ankyra/escape-inventory/metrics"
	"github.com/ankyra/escape-inventory/model"
	"github.com/gorilla/mux"
)

type uploadHandlerProvider struct {
	UploadPackage func(project, releaseId string, pkg io.ReadSeeker) error
}

func newUploadHandlerProvider() *uploadHandlerProvider {
	return &uploadHandlerProvider{
		UploadPackage: model.UploadPackage,
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	newUploadHandlerProvider().UploadHandler(w, r)
}

func (h *uploadHandlerProvider) UploadHandler(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	releaseId := name + "-" + version
	f, _, err := r.FormFile("file")
	if err != nil {
		HandleError(w, r, model.NewUserError(err))
		return
	}
	if err := h.UploadPackage(project, releaseId, f); err != nil {
		HandleError(w, r, err)
		return
	}
	metrics.UploadCounter.Inc()
	username := ReadUsernameFromContext(r)
	go CallWebHook(project, name, version, releaseId, username)
	w.WriteHeader(200)
}

func CallWebHook(project, unit, version, releaseId, username string) {
	if cmd.Config == nil || cmd.Config.WebHook == "" {
		return
	}
	prj := types.NewProject(project)
	app := types.NewApplication(project, unit)
	prjHooks, err := dao.GetProjectHooks(prj)
	if err != nil {
		log.Println("ERROR: Failed to get Inventory Project Hooks from database:", err)
		return
	}
	unitHooks, err := dao.GetApplicationHooks(app)
	if err != nil {
		log.Println("ERROR: Failed to get Inventory Application Hooks from database:", err)
		return
	}
	downstreamHooks, err := dao.GetDownstreamHooks(app)
	if err != nil {
		log.Println("ERROR: Failed to get Upstream Hooks from database:", err)
		return
	}
	fmt.Println("Downstream hooks:")
	fmt.Println(downstreamHooks)
	for _, hooks := range downstreamHooks {
		fmt.Println(hooks)
	}
	url := cmd.Config.WebHook
	data := map[string]interface{}{
		"event":            "NEW_UPLOAD",
		"project":          project,
		"project_hooks":    prjHooks,
		"unit":             unit,
		"unit_hooks":       unitHooks,
		"downstream_hooks": downstreamHooks,
		"version":          version,
		"release":          project + "/" + releaseId,
		"username":         username,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR: Failed to marshal webhook request:", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Println("ERROR: Failed to create webhook request:", err)
		return
	}
	log.Println("INFO: Calling webhook:", url)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: Failed to call webhook:", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Println("ERROR: Failed to call webhook, expecting status code 200, but got", resp.StatusCode)
		return
	}
}
