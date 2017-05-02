package handlers

import (
	"github.com/ankyra/escape-registry/model"
	"io/ioutil"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	metadata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	if err := model.AddRelease(string(metadata)); err != nil {
		HandleError(w, r, err)
		return
	}
}
