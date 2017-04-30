package handlers

import (
    "io/ioutil"
	"net/http"
    "encoding/json"
    "github.com/ankyra/escape-registry/model"
)

func ImportReleasesHandler(w http.ResponseWriter, r *http.Request) {
    releases, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    releasesList := []map[string]interface{}{}
    if err := json.Unmarshal(releases, &releasesList); err != nil {
        panic(err)
    }
    if err := model.Import(releasesList); err != nil {
        panic(err)
    }
}
