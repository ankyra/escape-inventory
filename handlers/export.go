package handlers

import (
	"net/http"
    "github.com/ankyra/escape-registry/model"
)

func ExportReleasesHandler(w http.ResponseWriter, r *http.Request) {
    if err := model.Export(w); err != nil {
        panic(err)
    }
}
