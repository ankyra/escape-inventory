package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/model"
)

func ErrorOrJsonSuccess(w http.ResponseWriter, r *http.Request, resp interface{}, err error) {
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JsonSuccess(w, resp)
}

func JsonSuccess(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

func ReadUsernameFromContext(r *http.Request) string {
	user := r.Context().Value("user")
	if user != nil {
		value := reflect.Indirect(reflect.ValueOf(user))
		return value.FieldByName("Name").String()
	}
	return ""
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		log.Println("Received nil error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if dao.IsNotFound(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if dao.IsAlreadyExists(err) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Resource already exists"))
		return
	} else if model.IsUserError(err) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println("Error:", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}
