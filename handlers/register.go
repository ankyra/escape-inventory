package handlers

import (
    "github.com/ankyra/escape-registry/model"
	"net/http"
    "io/ioutil"
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
