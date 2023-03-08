package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

func ValidateRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ServeHTTP(w, r)
		}
		secret := viper.GetViper().GetString("hmac_secret")
		//secret := "mysecret"

		byteData, err := io.ReadAll(r.Body)

		if err != nil {
			errMsg := "Bad data"
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		r.Body.Close() //  must close
		r.Body = io.NopCloser(bytes.NewBuffer(byteData))
		hmac := hmac.New(sha256.New, []byte(secret))

		hmac.Write(byteData)

		sha := hex.EncodeToString(hmac.Sum(nil))

		sig := r.Header.Get("Signature")

		if sha != sig {
			errMsg := "Forbidden"
			http.Error(w, errMsg, http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)

	})
}
