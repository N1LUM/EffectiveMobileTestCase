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

// UpdateUserByID обновляет данные пользователя по его ID.
//
// Swagger: operationId=updateUserById
// parameters:
//   - name: id
//     in: path
//     description: ID пользователя, чьи данные требуется обновить.
//     required: true
//     schema:
//     type: string
//     example: "550e8400-e29b-41d4-a716-446655440000"
//   - name: body
//     in: body
//     description: Новые данные пользователя для обновления.
//     required: true
//     schema:
//     "$ref": "#/definitions/UserUpdateInput"
//
// responses:
//
//	'200':
//	  description: Успешное обновление данных пользователя.
//	  schema:
//	    type: object
//	    properties:
//	      user_id:
//	        type: string
//	        example: "550e8400-e29b-41d4-a716-446655440000"
//	      msg:
//	        type: string
//	        example: "Обновление данных пользователя прошло успешно"
//	'400':
//	  description: Ошибка в запросе, например, неверный формат параметров или ошибка при обновлении данных.
//	  schema:
//	    type: object
//	    properties:
//	      error:
//	        type: string
//	        example: "Не удалось обновить данные пользователя: текст ошибки"
//
// definitions:
//
//	UserUpdateInput:
//	  type: object
//	  properties:
//	    Name:
//	      type: string
//	      example: "Иван"
//	    Surname:
//	      type: string
//	      example: "Иванов"
//	    Patronymic:
//	      type: string
//	      example: "Иванович"
//	    Address:
//	      type: string
//	      example: "г. Москва, ул. Пушкина, д. Колотушкина"
//	    PassportSerie:
//	      type: string
//	      example: "1234"
//	    PassportNumber:
//	      type: string
//	      example: "567890"
//	    CreatedAt:
//	      type: string
//	      format: date-time
//	      example: "2024-07-01T10:00:00Z"
//	    UpdatedAt:
//	      type: string
//	      format: date-time
//	      example: "2024-07-05T15:30:00Z"
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
