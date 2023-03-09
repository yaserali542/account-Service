package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/yaserali542/account-Service/controllers"
	"github.com/yaserali542/account-Service/repository"
	"github.com/yaserali542/account-Service/services"
)

func main() {

	var db *gorm.DB
	var err error
	v := initializeViperConfig()
	if db, err = repository.InitializeDatabase(v); err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	if err = repository.MigrateAccountTable(db); err != nil {
		log.Fatal(err)
		return
	}

	rep := repository.Repository{Db: db}
	service := services.AccountService{Repository: &rep}
	controller := controllers.Controllers{Services: service}
	middleware := controllers.Middleware{Services: service}

	r := mux.NewRouter()
	r.Use(middleware.ValidateRequest)
	r.HandleFunc("/refresh-token", controller.RefreshToken).Methods("GET")
	r.HandleFunc("/signin", controller.Signin).Methods("POST")
	r.HandleFunc("/signup", controller.SignUp).Methods("POST")

	r.HandleFunc("/testToken", middleware.ValidateJwtToken(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	log.Fatal(http.ListenAndServe(":8000", r))
}

func initializeViperConfig() *viper.Viper {
	viper.SetConfigType("json")
	viper.SetConfigFile("./config/config.json")
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	viper.ReadInConfig()
	return viper.GetViper()
}
