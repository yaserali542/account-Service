package services

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"

	"github.com/spf13/viper"
	"github.com/yaserali542/account-Service/models"
	"github.com/yaserali542/account-Service/repository"
	"golang.org/x/crypto/scrypt"
)

type AccountService struct {
	Repository *repository.Repository
}

func (service *AccountService) ValidateCredentials(creds models.Credentials) (*models.Account, error) {
	accountDetails, accountNotExist, err := service.Repository.GetUserDetails(creds.Username)

	if err != nil {
		return nil, err
	}
	if accountNotExist {
		return nil, errors.New("account not exists")
	}

	if accountDetails.UserName != creds.Username {
		return nil, errors.New("username mismatch")
	}

	hashPassword := generateHashedPassword(creds.Password, accountDetails.Salt)

	if !bytes.Equal(hashPassword, accountDetails.HashedPassword) {
		return nil, errors.New("invalid password")
	}
	return accountDetails, nil

}

func (service *AccountService) GetMinimalUserInfo(username string) (*models.BasicFields, error) {

	return service.Repository.GetMinimalUserInfo(username)

}

func (service *AccountService) CreateAccount(registerAccount *models.RegisterAccount) (*models.Account, bool, error) {
	accountDetails, accountNotExist, err := service.Repository.GetUserDetails(registerAccount.UserName)

	if !accountNotExist {
		return accountDetails, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	salt := generateSalt()

	hashedPassword := generateHashedPassword(registerAccount.Password, salt)

	account := &models.Account{
		FirstName:      registerAccount.FirstName,
		LastName:       registerAccount.LastName,
		UserName:       registerAccount.UserName,
		EmailAddress:   registerAccount.EmailAddress,
		ProfilePicture: registerAccount.ProfilePicture,
		HashedPassword: hashedPassword,
		Salt:           salt,
	}
	accountDetails, err = service.Repository.CreateAccount(account)
	if err != nil {
		return nil, true, err
	}
	return accountDetails, true, nil
}

func generateHashedPassword(password string, salt []byte) []byte {
	hashBytes := viper.GetViper().GetInt("pw_hash_bytes")
	hash, _ := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, hashBytes)
	return hash
}
func generateSalt() []byte {

	salt := make([]byte, viper.GetViper().GetInt("pw_salt_bytes"))
	io.ReadFull(rand.Reader, salt)

	return salt

}
