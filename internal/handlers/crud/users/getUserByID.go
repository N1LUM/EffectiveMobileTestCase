package users

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
)

// GetUserByID получает данные пользователя по его ID.
//
// Swagger: operationId=getUserByID
// parameters:
// - name: id
//   in: path
//   description: ID пользователя для получения данных.
//   required: true
//   type: string
// responses:
//   '200':
//     description: Успешное получение данных пользователя.
//     schema:
//       type: object
//       properties:
//         ID:
//           type: string
//           example: "550e8400-e29b-41d4-a716-446655440000"
//         Name:
//           type: string
//           example: "Иван"
//         Surname:
//           type: string
//           example: "Иванов"
//         Patronymic:
//           type: string
//           example: "Иванович"
//         Address:
//           type: string
//           example: "г. Москва, ул. Пушкина, д. Колотушкина"
//         PassportSerie:
//           type: string
//           example: "1234"
//         PassportNumber:
//           type: string
//           example: "567890"
//         FullPassport:
//           type: string
//           example: "1234567890"
//         CreatedAt:
//           type: string
//           format: date-time
//           example: "2024-07-05T12:00:00Z"
//         UpdatedAt:
//           type: string
//           format: date-time
//           example: "2024-07-05T12:30:00Z"
//   '400':
//     description: Ошибка в запросе, например, неверный формат ID пользователя или ошибка получения данных.
//     schema:
//       type: object
//       properties:
//         error:
//           type: string
//           example: "Не удалось получить данные пользователя: текст ошибки"

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на получение пользователя по ID")

	vars := mux.Vars(r)
	id := vars["id"]

	logging.Log.Debugf("ID пользователя на получение %v", id)

	var resultUser models.Users
	if err := db.PostgresClient.Where("ID = ?", id).Find(&resultUser).Error; err != nil {
		logging.Log.Errorf("Не удалось получить данные пользователя %v", err)
		http.Error(w, fmt.Sprintf("Не удалось получить данные пользователя: %v", err), 400)
		return
	}

	logging.Log.WithFields(logrus.Fields{
		"user_id":             resultUser.ID,
		"user_name":           resultUser.Name,
		"user_surname":        resultUser.Surname,
		"user_patronymic":     resultUser.Patronymic,
		"user_address":        resultUser.Address,
		"user_passportSerie":  resultUser.PassportSerie,
		"user_pussportNumber": resultUser.PassportNumber,
		"user_fullPassport":   resultUser.FullPassport,
		"user_createdAt":      resultUser.CreatedAt,
		"user_updatedAt":      resultUser.UpdatedAt,
	}).Debug("Получен пользователь со следующими данными")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"ID":             resultUser.ID.String(),
		"Name":           resultUser.Name,
		"Surname":        resultUser.Surname,
		"Patronymic":     resultUser.Patronymic,
		"Address":        resultUser.Address,
		"PassportSerie":  resultUser.PassportSerie,
		"PassportNumber": resultUser.PassportNumber,
		"FullPassport":   resultUser.FullPassport,
		"CreatedAt":      resultUser.CreatedAt.String(),
		"UpdatedAt":      resultUser.UpdatedAt.String(),
	})

	logging.Log.Info("Запрос на получение пользователя по ID успешно завершён")

	return
}
