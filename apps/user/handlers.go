package user

import (
	"errors"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/json"
	"go-base-structure/pkg/validators"
	"net/http"
	"strconv"
)

// Login log user in
func Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email" required:"true"`
		Password string `json:"password" required:"true"`
	}

	if validator := json.ReadJSON(w, r, &requestPayload); !validator.Valid() {
		if err := json.ErrorMapJSON(w, validator.Errors); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			userApp.Logger.Error("unable to write error json: ", err)
		}
		return
	}

	user, err := userApp.Models.User.GetUserByEmail(requestPayload.Email)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable get user from database: ", err)
		return
	}

	// validation for password goes here...
	validator := validators.New()
	isMatchPassword := validator.PasswordMatchesValidation(user.Password, requestPayload.Password)

	if !isMatchPassword {
		if err = json.ErrorStrJSON(w, errors.New("incorrect email or password")); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			userApp.Logger.Error("unable to write error json: ", err)
		}
		return
	}

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

	if err = json.WriteJSON(w, http.StatusOK, &tokens); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to write json: ", err)
		return
	}
}

// RefreshToken receive not expired refresh token and return new access token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		RefreshToken string `json:"refresh"`
	}

	var responseBody struct {
		AccessToken string `json:"access"`
	}

	if validator := json.ReadJSON(w, r, &requestPayload); !validator.Valid() {
		if err := json.ErrorMapJSON(w, validator.Errors); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			userApp.Logger.Error("unable to write error json: ", err)
		}
		return
	}

	claims, err := userApp.Auth.ParseWithClaims(requestPayload.RefreshToken)
	if err != nil || claims.TokenType != "refresh" {
		if err = json.ErrorStrJSON(w, errors.New("token is invalid or has expired"), http.StatusUnauthorized); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			userApp.Logger.Error("unable to write error json: ", err)
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
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("cannot get user from the database: ", err)
		return
	}

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
	if err = json.WriteJSON(w, http.StatusOK, &responseBody); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		userApp.Logger.Error("unable to write json: ", err)
		return
	}
}
