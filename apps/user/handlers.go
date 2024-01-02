package user

import (
	"errors"
	"go-base-structure/models"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/json"
	"go-base-structure/pkg/validators"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// Register do registration for users
func Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email           string `json:"email" required:"true"`
		FirstName       string `json:"first_name" required:"true"`
		LastName        string `json:"last_name" required:"true"`
		Password        string `json:"password" required:"true"`
		ConfirmPassword string `json:"confirm_password" required:"true"`
	}

	var responseBody struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	// get user data and json validation
	if validator := json.ReadJSON(w, r, &requestPayload); !validator.Valid() {
		if err := json.ErrorMapJSON(w, validator.Errors); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// validate email
	validator := validators.New()
	email := requestPayload.Email
	validator.EmailValidation(email)
	if !validator.Valid() {
		if err := json.ErrorMapJSON(w, validator.Errors); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// validate password
	validator = validators.New()
	password1 := requestPayload.Password
	password2 := requestPayload.ConfirmPassword

	validator.UserPasswordValidation(password1, password2)
	if !validator.Valid() {
		if err := json.ErrorMapJSON(w, validator.Errors); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// check if user with this information already exists
	isExistsUser, err := userApp.Models.User.CheckIfExistsUser(email)
	if err != nil {
		userApp.Logger.ServerError(w, "cannot check if user exists already", err)
		return
	}

	if isExistsUser {
		if err = json.ErrorStrJSON(w, errors.New("user with this email already exists")); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// insert user to the database
	user := &models.User{
		Email:     email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  password1,
	}

	_, err = userApp.Models.User.InsertOneUser(user)
	if err != nil {
		userApp.Logger.ServerError(w, "cannot insert user to the database", err)
		return
	}

	// write tokens to the output
	responseBody.Error = false
	responseBody.Message = "register was successful"
	if err = json.WriteJSON(w, http.StatusOK, responseBody); err != nil {
		userApp.Logger.ServerError(w, "unable to write json", err)
		return
	}
}

// Login log user in
func Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email" required:"true"`
		Password string `json:"password" required:"true"`
	}

	// read json
	if validator := json.ReadJSON(w, r, &requestPayload); !validator.Valid() {
		if err := json.ErrorMapJSON(w, validator.Errors); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// get user from the database
	user, err := userApp.Models.User.GetUserByEmail(requestPayload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = json.ErrorStrJSON(w, errors.New("incorrect email or password"), http.StatusUnauthorized); err != nil {
				userApp.Logger.ServerError(w, "unable to write error json", err)
			}
		} else {
			userApp.Logger.ServerError(w, "unable get user from database", err)
		}
		return
	}

	// validation for password goes here...
	validator := validators.New()
	validator.PasswordMatchesHashValidation(user.Password, requestPayload.Password)

	if !validator.Valid() {
		if err = json.ErrorStrJSON(w, errors.New("incorrect email or password"), http.StatusUnauthorized); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// generate new tokens
	jwtUser := &auth.JwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	tokens, err := userApp.Auth.GenerateTokenPair(jwtUser)
	if err != nil {
		userApp.Logger.ServerError(w, "unable to generate tokens", err)
		return
	}

	// write tokens to the output
	if err = json.WriteJSON(w, http.StatusOK, tokens); err != nil {
		userApp.Logger.ServerError(w, "unable to write json", err)
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

	// read json
	if validator := json.ReadJSON(w, r, &requestPayload); !validator.Valid() {
		if err := json.ErrorMapJSON(w, validator.Errors); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// validate token
	claims, err := userApp.Auth.ParseWithClaims(requestPayload.RefreshToken)
	if err != nil || claims.TokenType != "refresh" {
		if err = json.ErrorStrJSON(w, errors.New("token is invalid or has expired"), http.StatusUnauthorized); err != nil {
			userApp.Logger.ServerError(w, "unable to write error json", err)
		}
		return
	}

	// get user from the database
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		userApp.Logger.ServerError(w, "unable to convert string to int", err)
		return
	}

	user, err := userApp.Models.User.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = json.ErrorStrJSON(w, errors.New("incorrect email or password"), http.StatusUnauthorized); err != nil {
				userApp.Logger.ServerError(w, "unable to write error json", err)
			}
		} else {
			userApp.Logger.ServerError(w, "unable get user from database", err)
		}
		return
	}

	// generate new tokens
	jwtUser := &auth.JwtUser{
		ID:        userID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	tokens, err := userApp.Auth.GenerateTokenPair(jwtUser)
	if err != nil {
		userApp.Logger.ServerError(w, "unable to generate tokens", err)
		return
	}

	// write new refresh token to the output
	responseBody.AccessToken = tokens.Token
	if err = json.WriteJSON(w, http.StatusOK, responseBody); err != nil {
		userApp.Logger.ServerError(w, "unable to write json", err)
		return
	}
}
