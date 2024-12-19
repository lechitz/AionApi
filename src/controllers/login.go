package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/lechitz/AionApi/src/auth"
	"github.com/lechitz/AionApi/src/database"
	"github.com/lechitz/AionApi/src/models"
	"github.com/lechitz/AionApi/src/repository"
	"github.com/lechitz/AionApi/src/responses"
	"github.com/lechitz/AionApi/src/security"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User

	if err = json.Unmarshal(requestBody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositoryDB := repository.NewUserRepository(db)
	userDB, err := repositoryDB.GetUserByEmail(user.Email)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyPassword(userDB.Password, user.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.CreateToken(userDB.ID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(fmt.Sprintf("You are logged in: %s", token)))
}
