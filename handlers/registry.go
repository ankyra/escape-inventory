package handlers

import (
	"net/http"
    "encoding/json"
	"github.com/gorilla/mux"
    "github.com/ankyra/escape-registry/model"
)

func RegistryHandler(w http.ResponseWriter, r *http.Request) {
    typ := mux.Vars(r)["type"]
    name := mux.Vars(r)["name"]
    result, err := model.Registry(typ, name)
    if err != nil {
        HandleError(w, r, err)
        return
    }
    bytes, err := json.Marshal(result)
    if err != nil {
        HandleError(w, r, err)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write(bytes)
}
