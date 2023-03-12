package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/yaserali542/account-Service/models"
)

type Repository struct {
	Db *gorm.DB
}

func (repo *Repository) GetUserDetails(username string) (*models.Account, bool, error) {

	var account models.Account
	db := repo.Db.First(&account, "user_name = ?", username)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, true, nil
		} else {
			return nil, true, db.Error
		}

	}

	return &account, false, nil

}

func (repo *Repository) CreateAccount(account *models.Account) (*models.Account, error) {

	if err := repo.Db.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (repo *Repository) GetMinimalUserInfo(username string) (*models.BasicFields, error) {
	var account models.Account
	db := repo.Db.First(&account, "user_name = ?", username)
	if db.Error != nil {
		return nil, db.Error
	}
	return &models.BasicFields{
		ID:           account.ID,
		EmailAddress: account.EmailAddress,
		UserName:     account.UserName,
		Role:         account.Role,
	}, nil
}

func (repo *Repository) GetUserInfoFromId(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	db := repo.Db.First(&account, "id = ?", id)
	if db.Error != nil {
		return nil, db.Error
	}
	return &account, nil
}
