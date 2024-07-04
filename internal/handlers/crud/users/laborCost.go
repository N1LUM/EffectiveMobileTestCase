package users

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
	"time"
)

type TaskResponse struct {
	TaskID  uuid.UUID `json:"taskID"`
	Name    string    `json:"name"`
	Hours   int       `json:"hours"`
	Minutes int       `json:"minutes"`
	Seconds int       `json:"seconds"`
}

type Period struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type ByDuration []models.Tasks

func (a ByDuration) Len() int      { return len(a) }
func (a ByDuration) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool {
	if a[i].Hours != a[j].Hours {
		return a[i].Hours > a[j].Hours
	}
	if a[i].Minutes != a[j].Minutes {
		return a[i].Minutes > a[j].Minutes
	}
	return a[i].Seconds > a[j].Seconds
}

func LaborCost(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на получение трудозатрат пользователя")

	vars := mux.Vars(r)
	user_id := vars["user_id"]

	logging.Log.Debugf("ID пользователя, для которого будут получены трудозатраты %v", user_id)

	// Получаем параметры периода из тела запроса
	var period Period
	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		logging.Log.Errorf("Не удалось декодировать параметры периода: %v", err)
		http.Error(w, fmt.Sprintf("Не удалось декодировать параметры периода: %v", err), 400)
		return
	}

	logging.Log.Debugf("Период времени: %v - %v", period.StartTime, period.EndTime)

	//Получаем ID задач назначенных на пользователя
	var usersTasks []models.UsersTasks
	if err := db.PostgresClient.Where("user_id = ?", user_id).Find(&usersTasks).Error; err != nil {
		logging.Log.Errorf("Не удалось получить задачи пользователя %v", err)
		http.Error(w, fmt.Sprintf("Не удалось получить задачи пользователя: %v", err), 400)
		return
	}

	logging.Log.Debugf("Получен список отношений пользователя и задач: %v", usersTasks)

	//Получаем полные данные этих задач за период
	taskIDs := []uuid.UUID{}
	for _, taskRelation := range usersTasks {
		taskIDs = append(taskIDs, taskRelation.TaskID)
	}

	var tasks []models.Tasks
	if err := db.PostgresClient.Where("ID IN (?) AND start_time >= ? AND end_time <= ?", taskIDs, period.StartTime, period.EndTime).Find(&tasks).Error; err != nil {
		logging.Log.Errorf("Не удалось получить задачи: %v", err)
		http.Error(w, fmt.Sprintf("Не удалось получить задачи: %v", err), 400)
		return
	}

	logging.Log.Debugf("Получен список задач: %v", tasks)

	for index := range tasks {
		calculateTaskDuration(&tasks[index])
		logging.Log.WithFields(logrus.Fields{
			"TaskID":  tasks[index].ID,
			"Hours":   tasks[index].Hours,
			"Minutes": tasks[index].Minutes,
			"Seconds": tasks[index].Seconds,
		}).Debugf("Получена трудозатрата для пользователя с ID:  %v", user_id)
	}
	//Сортируем от большей к меньшей
	sort.Sort(ByDuration(tasks))

	// Создаем массив для ответа
	taskResponse := []TaskResponse{}
	for _, task := range tasks {
		taskResponse = append(taskResponse, TaskResponse{
			TaskID:  task.ID,
			Name:    task.Name,
			Hours:   task.Hours,
			Minutes: task.Minutes,
			Seconds: task.Seconds,
		})
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(taskResponse)

	logging.Log.Info("Запрос на получение трудозатрат пользователя успешно завершен")

}

func calculateTaskDuration(task *models.Tasks) {
	if task.StartTime != nil && task.EndTime != nil {
		duration := task.EndTime.Sub(*task.StartTime)
		totalSeconds := int(duration.Seconds())

		task.Hours = totalSeconds / 3600
		totalSeconds %= 3600

		task.Minutes = totalSeconds / 60
		task.Seconds = totalSeconds % 60
	}
}
