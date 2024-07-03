package crud

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
)

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
