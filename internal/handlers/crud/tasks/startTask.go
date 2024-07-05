package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
	"time"
)

// StartTaskTimer godoc
// @Summary      Start a task timer
// @Description  Start the timer for a specific task.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID"
// @Success      200  {object}  models.Tasks
// @Failure      400  {object}  map[string]string{"error": "Error message"}
// @Failure      404  {object}  map[string]string{"error": "Error message"}
// @Failure      500  {object}  map[string]string{"error": "Error message"}
// @Router       /tasks/{id}/start [post]
func StartTaskTimer(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на старт задачи")

	vars := mux.Vars(r)
	id := vars["id"]

	logging.Log.Debugf("ID задачи для старта %v", id)

	var task models.Tasks
	if err := db.PostgresClient.First(&task, "id = ?", id).Error; err != nil {
		logging.Log.Errorf("Задача не найдена: %v", err)
		http.Error(w, "Задача не найдена", 404)
		return
	}

	logging.Log.WithFields(logrus.Fields{
		"task_id":          task.ID,
		"task_name":        task.Name,
		"task_description": task.Description,
		"task_status":      task.Status,
		"task_hours":       task.Hours,
		"task_minutes":     task.Minutes,
		"task_seconds":     task.Seconds,
		"task_startTime":   task.StartTime,
		"task_endTime":     task.EndTime,
		"task_createdAt":   task.CreatedAt,
		"task_updatedAt":   task.UpdatedAt,
	}).Debug("Была найдена задача")

	startTime := time.Now()
	task.StartTime = &startTime
	task.Status = true

	if err := db.PostgresClient.Save(&task).Error; err != nil {
		logging.Log.Errorf("Ошибка при обновлении задачи: %v", err)
		http.Error(w, fmt.Sprintf("Ошибка при обновлении задачи: %v", err), 500)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(task)

	logging.Log.Infof("Отсчет времени для задачи %s начат", id)

	logging.Log.Info("Запрос на старт задачи успешно завершен")

	return
}
