package services

import (
	"bytes"
	"errors"

	"github.com/spf13/viper"
	"github.com/yaserali542/account-Service/models"
	"github.com/yaserali542/account-Service/repository"
	"golang.org/x/crypto/scrypt"
)

type AccountService struct {
	Repository *repository.Repository
}

func (service *AccountService) ValidateCredentials(creds models.Credentials) (*models.Account, error) {
	accountDetails, err := service.Repository.GetUserDetails(creds.Username)

	if err != nil {
		return nil, err
	}

	if accountDetails.UserName != creds.Username {
		return nil, errors.New("username mismatch")
	}

	hashPassword := generateHashWithOldSalt(creds.Password, accountDetails.Salt)

	if !bytes.Equal(hashPassword, accountDetails.HashedPassword) {
		return nil, errors.New("invalid password")
	}
	return accountDetails, nil

}

func generateHashWithOldSalt(password string, salt []byte) []byte {
	hashBytes := viper.GetViper().GetInt("pw_hash_bytes")
	hash, _ := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, hashBytes)
	return hash

}
