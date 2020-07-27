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
	Date          time.Time     `validate:"required"`
	OperationType OperationType `validate:"required,numeric,gte=1,lte=4"`
	Amount        float64       `validate:"required,number"`
	Remark        string
	AccountID     uuid.UUID `gorm:"type:varchar(255);" validate:"required"`
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
	return [...]string{"", "Deposit", "Withdrawal", "Stocks Purchase", "Stocks Sales"}[ot]
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

// FundFlowCreate : Create a new fund flow transaction for the user account
func FundFlowCreate(fundFlow *FundFlow) error {
	// Validate input first
	err := validate.Struct(fundFlow)

	if err != nil {
		if err = getValidationMessage(err); err != nil {
			return err
		}
	}

	db := GetDB()
	defer db.Close()

	db.Save(fundFlow)

	if fundFlow.ID == uuid.Nil {
		return fmt.Errorf("fund flow is not saved")
	}

	fundFlow.Operation = fmt.Sprintf("%v", OperationType(fundFlow.OperationType))

	return nil
}

// FundFlowUpdate : Update an existing fund flow
func FundFlowUpdate(fundFlow *FundFlow) error {
	// Validate input first
	err := validate.Struct(fundFlow)

	if err != nil {
		if err = getValidationMessage(err); err != nil {
			return err
		}
	}

	db := GetDB()
	defer db.Close()

	temp := FundFlow{}
	db.Table("fund_flows").Where("id = ?", fundFlow.ID).First(&temp)

	if temp.ID == uuid.Nil {
		return fmt.Errorf("fund flow not found")
	}

	temp.Date = fundFlow.Date
	temp.OperationType = fundFlow.OperationType
	temp.Amount = fundFlow.Amount
	temp.Remark = fundFlow.Remark
	db.Save(&temp)

	fundFlow.Operation = fmt.Sprintf("%v", OperationType(temp.OperationType))

	return nil
}

// FundFlowDelete : Delete an existing fund flow
func FundFlowDelete(fundFlow *FundFlow) error {
	db := GetDB()
	defer db.Close()

	temp := FundFlow{}
	db.Table("fund_flows").Where("id = ?", fundFlow.ID).First(&temp)

	if temp.ID == uuid.Nil {
		return fmt.Errorf("fund flow not found")
	}

	db.Delete(&temp)

	return nil
}
