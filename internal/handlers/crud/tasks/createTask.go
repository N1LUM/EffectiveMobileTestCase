package tasks

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
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на создание задания")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось прочитать тело запроса")
		http.Error(w, fmt.Sprintf("Не удалось прочитать тело запроса: %v", err), 400)
		return
	}

	var task models.Tasks
	task.ID, err = uuid.NewUUID()
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось сгенерировать uuid для нового задания")
		http.Error(w, fmt.Sprintf("Не удалось сгенерировать uuid для нового задания: %v", err), 500)
		return
	}

	logging.Log.Debugf("Сгенерирован uuid для нового задания - %v", task.ID)

	if err = json.Unmarshal(body, &task); err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось декодировать тело запроса в структуру Tasks")
		http.Error(w, fmt.Sprintf("Не удалось декодировать тело запроса в структуру Tasks: %v", err), 400)
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
	}).Debug("Данные для создания записи задания")

	if err = db.PostgresClient.Create(&task).Error; err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Неудалось создать задание")
		http.Error(w, fmt.Sprintf("Неудалось создать задание: %v", err), 500)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID.String(), "msg": "Создание задания прошло успешно"})

	logging.Log.Info("Запрос на создание задания успешно завершен")

	return
}
