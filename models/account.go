package models

import (
	"fmt"
	"github.com/satori/go.uuid"
)

// Account : The fund account that the user has
type Account struct {
	Base
	Title  string
	UserID uuid.UUID `gorm:"type:varchar(255);"`
}

// AccountList : Get all the accounts that belong to the user
func AccountList(userID string) ([]Account, error) {
	db := GetDB()
	defer db.Close()

	accounts := []Account{}
	db.Table("accounts").Where("user_id = ?", userID).Find(&accounts)

	return accounts, nil
}

// AccountCreate : Create a new account for the user
func AccountCreate(account *Account) error {
	db := GetDB()
	defer db.Close()

	db.Save(account)

	if account.ID == uuid.Nil {
		return fmt.Errorf("account is not saved")
	}

	return nil
}

// AccountUpdate : Update an existing account
func AccountUpdate(account *Account) error {
	db := GetDB()
	defer db.Close()

	temp := Account{}
	db.Table("accounts").Where("id = ? and user_id = ?", account.ID, account.UserID).First(&temp)

	if temp.ID == uuid.Nil {
		return fmt.Errorf("account not found")
	}

	temp.Title = account.Title
	db.Save(&temp)

	return nil
}

// AccountDelete : Delete an existing account
func AccountDelete(account *Account) error {
	db := GetDB()
	defer db.Close()

	temp := Account{}
	db.Table("accounts").Where("id = ? and user_id = ?", account.ID, account.UserID).First(&temp)

	if temp.ID == uuid.Nil {
		return fmt.Errorf("account not found")
	}

	db.Delete(&temp)

	return nil
}
