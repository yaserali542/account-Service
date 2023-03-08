package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/yaserali542/account-Service/models"
	"github.com/yaserali542/account-Service/services"
)

type Controllers struct {
	Services services.AccountService
}

func (c *Controllers) Signin(w http.ResponseWriter, r *http.Request) {

	var creds models.Credentials
	// Get the JSON body and decode into credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}

	account, err := c.Services.ValidateCredentials(creds)

	if err != nil {
		errMsg := "Forbidden"
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(account)

}
