package users

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
	"test/internal/validation"
)

// UpdateUserByID godoc
// @Summary      Update user data by ID
// @Description  Update user data specified by the user ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        body body models.Users true "User object to update"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string{"error": "Error message"}
// @Router       /users/{id} [put]
func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на обновление данных пользователя")

	vars := mux.Vars(r)
	id := vars["id"]

	logging.Log.Debugf("ID пользователя на обновление данных %v", id)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось прочитать тело запроса")
		http.Error(w, fmt.Sprintf("Не удалось прочитать тело запроса: %v", err), 400)
		return
	}

	var user models.Users
	if err = json.Unmarshal(body, &user); err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось декодировать тело запроса в структуру Users")
		http.Error(w, fmt.Sprintf("Не удалось удалось декодировать тело запроса в структуру Users: %v", err), 400)
		return
	}

	user.FullPassport = user.PassportSerie + user.PassportNumber

	if err = validation.ValidateUpdateUser(&user); err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Данные пользователя не прошли валидацию")
		http.Error(w, fmt.Sprintf("Данные пользователя не прошли валидацию: %v", err), 400)
		return
	}

	logging.Log.WithFields(logrus.Fields{
		"user_id":             user.ID,
		"user_name":           user.Name,
		"user_surname":        user.Surname,
		"user_patronymic":     user.Patronymic,
		"user_address":        user.Address,
		"user_passportSerie":  user.PassportSerie,
		"user_pussportNumber": user.PassportNumber,
		"user_fullPassport":   user.FullPassport,
		"user_createdAt":      user.CreatedAt,
		"user_updatedAt":      user.UpdatedAt,
	}).Debugf("Данные для обновления записи пользователя с ID: %v", id)

	if err = db.PostgresClient.Where("ID = ?", id).Updates(user).Error; err != nil {
		logging.Log.Errorf("Не удалось обновить данные пользователя %v", err)
		http.Error(w, fmt.Sprintf("Не удалось обновить данные пользователя: %v", err), 400)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"user_id": id, "msg": "Обновление данных пользователя прошло успешно"})

	logging.Log.Info("Запрос на обновление данных пользователя успешно завершен")

	return
}
