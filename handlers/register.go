package handlers

import (
	"net/http"
    "io/ioutil"
    "github.com/ankyra/escape-registry/model"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    metadata, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    if err := model.AddRelease(string(metadata)); err != nil {
        panic(err)
    }
}
