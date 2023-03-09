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

func (c *Controllers) SignUp(w http.ResponseWriter, r *http.Request) {
	var registerAccount models.RegisterAccount

	if err := json.NewDecoder(r.Body).Decode(&registerAccount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdAccount, accountNotExist, err := c.Services.CreateAccount(&registerAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !accountNotExist {
		errMsg := "Account Exists"
		http.Error(w, errMsg, http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(createdAccount)

}
