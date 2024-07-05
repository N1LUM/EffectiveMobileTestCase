package users

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
)

// DeleteUserByID удаляет пользователя по его ID.
//
// Swagger: operationId=deleteUserByID
// parameters:
// - name: id
//   in: path
//   description: ID пользователя для удаления.
//   required: true
//   type: string
// responses:
//   '200':
//     description: Успешное удаление пользователя.
//     schema:
//       type: object
//       properties:
//         user_id:
//           type: string
//           example: "550e8400-e29b-41d4-a716-446655440000"
//         msg:
//           type: string
//           example: "Удаление пользователя прошло успешно"
//   '400':
//     description: Ошибка в запросе, например, неверный формат ID пользователя или ошибка удаления.
//     schema:
//       type: object
//       properties:
//         error:
//           type: string
//           example: "Не удалось удалить пользователя: текст ошибки"

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на удаление пользователя по ID")

	vars := mux.Vars(r)
	id := vars["id"]

	logging.Log.Debugf("ID пользователя на удаление %v", id)

	if err := db.PostgresClient.Where("ID = ?", id).Delete(&models.Users{}).Error; err != nil {
		logging.Log.Errorf("Не удалось удалить пользователя %v", err)
		http.Error(w, fmt.Sprintf("Не удалось удалить пользователя: %v", err), 400)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"user_id": id, "msg": "Удаление пользователя прошло успешно"})

	logging.Log.Info("Запрос на удаление пользователя по ID успешно завершён")

	return
}
