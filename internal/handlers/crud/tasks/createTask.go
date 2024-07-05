package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
)

// CreateTask создает новое задание для указанного пользователя.
//
// Swagger: operationId=createTask
// responses:
//   '200':
//     description: Успешное создание задания.
//     schema:
//       type: object
//       properties:
//         task_id:
//           type: string
//           example: "550e8400-e29b-41d4-a716-446655440000"
//         msg:
//           type: string
//           example: "Создание задания прошло успешно"
//   '400':
//     description: Ошибка в запросе, например, неверный формат UUID пользователя.
//   '500':
//     description: Внутренняя ошибка сервера при создании задания или связи пользователя и задания.
//     schema:
//       type: object
//       properties:
//         error:
//           type: string
//           example: "Не удалось создать задание: ошибка текст"

func CreateTask(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на создание задания")

	vars := mux.Vars(r)
	user_id, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось перевести ID пользователя с UUID в String")
		http.Error(w, fmt.Sprintf("Не удалось перевести ID пользователя с UUID в String: %v", err), 400)
		return
	}

	logging.Log.Debugf("Задача будет создана для пользователя с ID - %v", user_id)

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

	user_tasks := models.UsersTasks{
		TaskID: task.ID,
		UserID: user_id,
	}
	user_tasks.ID, err = uuid.NewUUID()
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Не удалось сгенерировать uuid для записи связи пользователя и задания")
		http.Error(w, fmt.Sprintf("Не удалось сгенерировать uuid для записи связи пользователя и задания: %v", err), 500)
		return
	}

	if err = db.PostgresClient.Create(&user_tasks).Error; err != nil {
		logging.Log.WithFields(logrus.Fields{
			"errors": err,
		}).Error("Неудалось создать связь между пользователем и заданием")
		http.Error(w, fmt.Sprintf("Неудалось создать связь между пользователем и заданием: %v", err), 500)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID.String(), "msg": "Создание задания прошло успешно"})

	logging.Log.Info("Запрос на создание задания успешно завершен")

	return
}
