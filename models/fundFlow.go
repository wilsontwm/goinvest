package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

// FundFlow : The fund flow that records all the money transaction that the user has performed
type FundFlow struct {
	Base
	Date          time.Time
	OperationType OperationType
	Amount        float64
	Remark        string
	AccountID     uuid.UUID `gorm:"type:varchar(255);"`
	Operation     string    `gorm:"-"`
}

// FundFlowListFilter : The filter used for fund flow list
type FundFlowListFilter struct {
	UserID     uuid.UUID
	AccountIDs []uuid.UUID
	Limit      int
	Page       int
}

// OperationType : The operation type (deposit, withdraw, stocks buy, stocks sell)
type OperationType int

const (
	// Deposit operation
	Deposit OperationType = iota + 1
	// Withdrawal operation
	Withdrawal
	// StocksPurchase operation
	StocksPurchase
	// StocksSales operation
	StocksSales
)

func (ot OperationType) String() string {
	return [...]string{"Deposit", "Withdrawal", "Stocks Purchase", "Stocks Sales"}[ot]
}

// FundFlowList : Get all the money transactions that belong to the user
func FundFlowList(filter FundFlowListFilter) ([]FundFlow, error) {
	db := GetDB()
	defer db.Close()

	fundFlows := []FundFlow{}
	// Get the list of fund flows
	db.Table("fund_flows").
		Joins("left join accounts on fund_flows.account_id = accounts.id").
		Scopes(filterByUserID(filter.UserID), filterByAccountIDs(filter.AccountIDs)).
		Order("date desc").
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&fundFlows)

	for i := range fundFlows {
		fundFlows[i].Operation = fmt.Sprintf("%v", OperationType(fundFlows[i].OperationType))
	}

	return fundFlows, nil
}

// Filter the fund flow by the user
func filterByUserID(userID uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		return db.Where("user_id = ?", userID)
	}
}

// Filter the fund flow by the account
func filterByAccountIDs(accountIDs []uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(accountIDs) > 0 {
			return db.Where("account_id IN (?)", accountIDs)
		}
		return db
	}
}

/*
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
*/
