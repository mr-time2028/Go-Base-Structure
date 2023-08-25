package user

import (
	"errors"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/json"
	"net/http"
	"strconv"
	"time"
)

// Login log user in
func Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string
		Password string
	}

	err := json.ReadJSON(w, r, &requestPayload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to read json: ", err)
		return
	}

	// validation for email goes here...

	user, err := userApp.Models.User.GetUserByEmail(requestPayload.Email)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable get user from database: ", err)
		return
	}

	// validation for password goes here...

	jwtUser := &auth.JwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	tokens, err := userApp.Auth.GenerateTokenPair(jwtUser)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to generate tokens: ", err)
		return
	}

	err = json.WriteJSON(w, http.StatusOK, &tokens)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to write json: ", err)
		return
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		RefreshToken string `json:"refresh_token"`
	}

	var responseBody struct {
		AccessToken string `json:"access_token"`
	}

	err := json.ReadJSON(w, r, &requestPayload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to read json: ", err)
		return
	}

	claims, err := userApp.Auth.ParseWithClaims(requestPayload.RefreshToken)
	if err != nil || claims.Issuer != userApp.Auth.Issuer || claims.ExpiresAt.Before(time.Now()) {
		if err = json.ErrorJSON(w, errors.New("token is invalid or has expired"), http.StatusUnauthorized); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			userApp.Logger.Error("unable to write error json: ", err)
			return
		}
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to convert string to int: ", err)
		return
	}

	user, err := userApp.Models.User.GetUserByID(userID)
	jwtUser := &auth.JwtUser{
		ID:        userID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	tokens, err := userApp.Auth.GenerateTokenPair(jwtUser)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to generate tokens: ", err)
		return
	}

	responseBody.AccessToken = tokens.Token
	err = json.WriteJSON(w, http.StatusOK, &responseBody)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to write json: ", err)
		return
	}
}
