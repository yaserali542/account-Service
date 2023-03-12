package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
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
	jwtToken := services.GenerateToken(account.UserName)

	jwtAccount := models.AccountWithJWT{
		FirstName:      account.FirstName,
		LastName:       account.LastName,
		UserName:       account.UserName,
		EmailAddress:   account.EmailAddress,
		Role:           account.Role,
		ProfilePicture: account.ProfilePicture,
		JwtToken: models.JwtToken{
			Token: jwtToken,
		},
	}

	if err != nil {
		errMsg := "Forbidden"
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(jwtAccount)

}

func (c *Controllers) SignUp(w http.ResponseWriter, r *http.Request) {
	var registerAccount models.RegisterAccount

	if err := json.NewDecoder(r.Body).Decode(&registerAccount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (*Controllers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	v := viper.GetViper()

	jwtTokenString := r.Header.Get("token")

	if jwtTokenString == "" {
		http.Error(w, "token is missing", http.StatusBadRequest)
		return
	}
	claims := &models.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(jwtTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		jwtKey := v.GetString("jwt.key")
		return []byte(jwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtToken := models.JwtToken{
		Token: services.GenerateToken(claims.Username),
	}
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(&jwtToken)

}

func (c *Controllers) GetMinimalUserInfo(w http.ResponseWriter, r *http.Request) {
	var username models.BasicFieldsRequest

	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	basicDetails, err := c.Services.GetMinimalUserInfo(username.UserName)
	if err != nil {
		errMsg := "internal server error"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(basicDetails)

}

func (c *Controllers) GetUserInfoById(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("this method is invoked")
	vars := mux.Vars(r)
	id, _ := uuid.FromString(vars["id"])

	accountDetails, err := c.Services.GetUserInfoFromId(id)
	if err != nil {
		errMsg := "internal server error"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(accountDetails)

}

func (c *Controllers) VerfierSignin(w http.ResponseWriter, r *http.Request) {

	var creds models.Credentials
	// Get the JSON body and decode into credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		return
	}

	account, err := c.Services.ValidateVerifierCredentials(creds)
	jwtToken := services.GenerateToken(account.UserName)

	jwtAccount := models.AccountWithJWT{
		FirstName:    account.FirstName,
		LastName:     account.LastName,
		UserName:     account.UserName,
		EmailAddress: account.EmailAddress,
		Role:         account.Role,
		//ProfilePicture: account.ProfilePicture,
		JwtToken: models.JwtToken{
			Token: jwtToken,
		},
	}

	if err != nil {
		errMsg := "Forbidden"
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(jwtAccount)

}
