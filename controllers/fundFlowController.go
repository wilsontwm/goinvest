package controllers

import (
	//"encoding/json"
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
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, funds, "Successfully retrieved the fund flows list.")

}

/*
// AccountCreate (Post): Create a new account
var AccountCreate = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)

	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	account.UserID = userID

	if err = models.AccountCreate(account); err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, &account, "You have successfully created a new account.")
}

// AccountUpdate (Patch): Update an existing account
var AccountUpdate = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)

	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	account.UserID = userID

	if err = models.AccountUpdate(account); err != nil {
		utils.Fail(w, http.StatusNotFound, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, &account, "You have successfully updated an existing account.")
}

// AccountDelete (Delete): Delete an existing account
var AccountDelete = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	data := make(map[string]interface{})
	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)

	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	account.UserID = userID

	if err = models.AccountDelete(account); err != nil {
		utils.Fail(w, http.StatusNotFound, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, data, "You have successfully deleted an existing account.")
}
*/
