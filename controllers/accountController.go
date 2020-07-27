package controllers

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"goinvest/models"
	"goinvest/utils"
	"net/http"
)

// AccountList (Get): Get all the accounts belonging to the user
var AccountList = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	userID := r.Context().Value(models.ContextKeyUserID).(uuid.UUID)

	accounts, err := models.AccountList(userID.String())

	if err != nil {
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, accounts, "Successfully retrieved the accounts.")

}

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
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, &account, "You have successfully created a new account.")
}

// AccountUpdate (Post): Update an existing account
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
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, &account, "You have successfully updated an existing account.")
}

// AccountDelete (Post): Delete an existing account
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
		utils.Fail(w, http.StatusUnprocessableEntity, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, data, "You have successfully deleted an existing account.")
}
