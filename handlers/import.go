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
	"encoding/json"
	"fmt"
	"github.com/ankyra/escape-inventory/model"
	"io/ioutil"
	"net/http"
)

func ImportReleasesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		HandleError(w, r, model.NewUserError(fmt.Errorf("Empty body")))
		return
	}
	releases, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	releasesList := []map[string]interface{}{}
	if err := json.Unmarshal(releases, &releasesList); err != nil {
		HandleError(w, r, model.NewUserError(err))
		return
	}
	if err := model.Import(releasesList); err != nil {
		HandleError(w, r, err)
		return
	}
}
