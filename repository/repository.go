package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/yaserali542/account-Service/models"
)

type Repository struct {
	Db *gorm.DB
}

func (repo *Repository) GetUserDetails(username string) (*models.Account, bool, error) {

	var account models.Account
	db := repo.Db.First(&account, "user_name = ?", username)
	if db.Error != nil && errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, true, nil
	}
	return &account, false, nil

}

func (repo *Repository) CreateAccount(account *models.Account) (*models.Account, error) {

	if err := repo.Db.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}
