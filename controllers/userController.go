package controllers

import (
	"encoding/json"
	"goinvest/models"
	"goinvest/utils"
	"net/http"
)

// UserLogin (POST): Login the user
var UserLogin = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	if err = models.UserLogin(user); err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, &user, "You have successfully login.")
}
