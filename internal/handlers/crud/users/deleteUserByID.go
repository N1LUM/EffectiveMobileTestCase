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

// DeleteUserByID godoc
// @Summary      Delete a user by ID
// @Description  Delete a user by ID
// @Tags         users
// @Param        id   path      string  true  "User ID"
// @Produce      json
// @Success      200  {object}  map[string]string{"user_id": "ID of the deleted user", "msg": "Success message"}
// @Failure      400  {object}  map[string]string{"error": "Error message"}
// @Router       /users/delete/{id} [delete]
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
