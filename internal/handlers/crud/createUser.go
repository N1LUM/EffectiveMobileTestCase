package crud

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
	"test/internal/validation"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на создание пользователя")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось прочитать тело запроса")
		http.Error(w, fmt.Sprintf("Не удалось прочитать тело запроса: %v", err), 400)
		return
	}

	var user models.Users

	user.ID, err = uuid.NewUUID()
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось сгенерировать uuid для нового пользователя")
		http.Error(w, fmt.Sprintf("Не удалось сгенерировать uuid для нового пользователя: %v", err), 500)
		return
	}

	logging.Log.Debugf("Сгенерирован uuid для нового пользователя")

	if err = json.Unmarshal(body, &user); err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось перевести байты тела запроса в структуру Users")
		http.Error(w, fmt.Sprintf("Не удалось перевести байты тела запроса в структуру Users: %v", err), 400)
		return
	}

	user.FullPassport = user.PassportSerie + user.PassportNumber

	if err = validation.ValidateUser(&user); err != nil {
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
	}).Debug("Данные для создания записи пользователя")

	if err = db.PostgresClient.Create(&user).Error; err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Неудалось создать пользователя")
		http.Error(w, fmt.Sprintf("Неудалось создать пользователя: %v", err), 500)
		return
	}
	logging.Log.Info("Запрос на создание пользователя успешно завершен")

	w.WriteHeader(http.StatusOK)

	return
}
