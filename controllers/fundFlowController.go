package controllers

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"goinvest/models"
	"goinvest/utils"
	"net/http"
	"strconv"
	"strings"
)

// FundFlowList (Get): Get all the fund flows belonging to the user
var FundFlowList = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	accountIDs := []uuid.UUID{}
	if accountIDsParam := utils.GetParam(r, "accountIDs"); accountIDsParam != "" {
		accountIDstrs := strings.Split(accountIDsParam, ",")
		for _, accountIDstr := range accountIDstrs {
			if accountID, err := uuid.FromString(accountIDstr); err == nil {
				accountIDs = append(accountIDs, accountID)
			}
		}
	}

	limit := 10
	if limitParam := utils.GetParam(r, "limit"); limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}

	page := 1
	if pageParam := utils.GetParam(r, "page"); pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	// Build the filter for the fund flow list
	filter := models.FundFlowListFilter{
		UserID:     userID,
		AccountIDs: accountIDs,
		Limit:      limit,
		Page:       page,
	}

	// Get the funds flow list
	funds, err := models.FundFlowList(filter)

	if err != nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, funds, "Successfully retrieved the fund flows list.")

}

// FundFlowCreate (Post): Create a new fund flow transaction
var FundFlowCreate = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	fundFlow := &models.FundFlow{}
	err := json.NewDecoder(r.Body).Decode(fundFlow)

	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	// Check if the account belongs to the user
	account, _ := models.AccountGet(userID, fundFlow.AccountID)

	if account.ID == uuid.Nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, "Invalid account ID.")
		return
	}

	if err = models.FundFlowCreate(fundFlow); err != nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, fundFlow, "You have successfully created a new fund flow transaction.")
}

// FundFlowUpdate (Post): Update an existing fund flow transaction
var FundFlowUpdate = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	fundFlow := &models.FundFlow{}
	err := json.NewDecoder(r.Body).Decode(fundFlow)

	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	// Check if the account belongs to the user
	account, _ := models.AccountGet(userID, fundFlow.AccountID)
	if account.ID == uuid.Nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, "Invalid account ID.")
		return
	}

	if err = models.FundFlowUpdate(fundFlow); err != nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, fundFlow, "You have successfully updated an existing fund flow transaction.")
}

// FundFlowDelete (Post): Delete an existing fund flow transaction
var FundFlowDelete = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	data := make(map[string]interface{})
	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	fundFlow := &models.FundFlow{}
	err := json.NewDecoder(r.Body).Decode(fundFlow)

	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	// Check if the account belongs to the user
	account, _ := models.AccountGet(userID, fundFlow.AccountID)
	if account.ID == uuid.Nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, "Invalid account ID.")
		return
	}

	if err = models.FundFlowDelete(fundFlow); err != nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, data, "You have successfully deleted an existing fund flow transaction.")
}
