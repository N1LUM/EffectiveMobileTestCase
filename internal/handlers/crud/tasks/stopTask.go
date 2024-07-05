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

// StopTaskTimer останавливает таймер для указанной задачи.
//
// Swagger: operationId=stopTaskTimer
// parameters:
// - name: id
//   in: path
//   description: ID задачи для остановки таймера.
//   required: true
//   type: string
// responses:
//   '200':
//     description: Успешная остановка таймера для задачи.
//     schema:
//       type: object
//       properties:
//         ID:
//           type: string
//           example: "550e8400-e29b-41d4-a716-446655440000"
//         Name:
//           type: string
//           example: "Название задачи"
//         Description:
//           type: string
//           example: "Описание задачи"
//         Status:
//           type: boolean
//           example: false
//         Hours:
//           type: integer
//           example: 0
//         Minutes:
//           type: integer
//           example: 0
//         Seconds:
//           type: integer
//           example: 0
//         StartTime:
//           type: string
//           format: date-time
//         EndTime:
//           type: string
//           format: date-time
//         CreatedAt:
//           type: string
//           format: date-time
//         UpdatedAt:
//           type: string
//           format: date-time
//   '404':
//     description: Задача с указанным ID не найдена.
//   '500':
//     description: Внутренняя ошибка сервера при обновлении задачи.
//     schema:
//       type: object
//       properties:
//         error:
//           type: string
//           example: "Ошибка при обновлении задачи: текст ошибки"

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
