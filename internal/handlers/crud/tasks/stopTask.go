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

func StopTaskTimer(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на остановку задачи")

	vars := mux.Vars(r)
	id := vars["id"]

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

	endTime := time.Now()
	task.EndTime = &endTime
	task.Status = false

	if err := db.PostgresClient.Save(&task).Error; err != nil {
		logging.Log.Errorf("Ошибка при обновлении задачи: %v", err)
		http.Error(w, fmt.Sprintf("Ошибка при обновлении задачи: %v", err), 500)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(task)

	logging.Log.Infof("Отсчет времени для задачи %s закончен", id)

	logging.Log.Info("Запрос на остановку задачи успешно завершен")

	return
}
