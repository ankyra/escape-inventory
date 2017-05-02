package handlers

import (
	"github.com/ankyra/escape-registry/model"
	"net/http"
)

func ExportReleasesHandler(w http.ResponseWriter, r *http.Request) {
	if err := model.Export(w); err != nil {
		HandleError(w, r, err)
		return
	}
}
