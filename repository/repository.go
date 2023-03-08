package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/yaserali542/account-Service/models"
)

type Repository struct {
	Db *gorm.DB
}

func (repo *Repository) GetUserDetails(username string) (*models.Account, error) {

	var account models.Account

	if err := repo.Db.First(&account, "user_name = ?", username).Error; err != nil {
		return nil, err
	}

	return &account, nil

}
