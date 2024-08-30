package api

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go-notes/models"
	"go-notes/services"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.UserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithError(err).Error("Invalid request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if password, ok := models.Users[user.Username]; !ok || password != user.Password {
		logrus.Warn("Invalid credentials for user: ", user.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Username,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		logrus.WithError(err).Error("Could not generate token")
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	logrus.Info("User logged in: ", user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func AddNote(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	var note models.UserNoteReq

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		logrus.WithError(err).Error("Invalid request payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка орфографии через Яндекс.Спеллер
	if err := services.CheckSpelling(note.Content); err != nil {
		logrus.WithError(err).Error("Spelling error in note")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	models.Notes[user] = append(models.Notes[user], note.Content)
	logrus.Info("Note added for user: ", user)
	w.WriteHeader(http.StatusCreated)
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)
	logrus.Info("Fetching notes for user: ", user)
	json.NewEncoder(w).Encode(models.Notes[user])
}
