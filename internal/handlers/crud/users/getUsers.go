package users

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
)

// GetUsers godoc
// @Summary      Get a list of users with optional filtering
// @Description  Get a list of users with optional filtering. Supports pagination.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        filters body models.UserGetListInput true "Filters to apply"
// @Success      200  {array}  models.Users
// @Failure      400  {object}  map[string]string{"error": "Error message"}
// @Router       /users/list [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Запрос на получение списка пользователей")

	// Извлекаем параметры фильтрации из запроса, если они есть
	input := models.UserGetListInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err != io.EOF {
		logging.Log.Errorf("Ошибка при декодировании параметров фильтрации: %v", err)
		http.Error(w, fmt.Sprintf("Ошибка при декодировании параметров фильтрации: %v", err), 400)
		return
	}

	logging.Log.WithFields(logrus.Fields{
		"page":    input.Page,
		"limit":   input.Limit,
		"filters": input.Filters,
	}).Debug("Получен Input со следующими данными")

	page := 1
	limit := 10

	if input.Page > 0 {
		page = input.Page
	}

	if input.Limit > 0 {
		limit = input.Limit
	}

	offset := (page - 1) * limit

	logging.Log.Debugf("Page=%d, Limit=%d, Offset=%d", page, limit, offset)

	query := db.PostgresClient.Model(&models.Users{})

	if len(input.Filters.Filters) > 0 {
		for _, filter := range input.Filters.Filters {
			switch filter.Operator {
			case "equals":
				if filter.Value != "" {
					query = applyFilter(query, filter.Field, filter.Value, filter.Operator)
				}
			case "contains":
				if filter.Value != "" {
					query = applyFilter(query, filter.Field, filter.Value, filter.Operator)
				}
			case "startsWith":
				if filter.Value != "" {
					query = applyFilter(query, filter.Field, filter.Value, filter.Operator)
				}
			case "endsWith":
				if filter.Value != "" {
					query = applyFilter(query, filter.Field, filter.Value, filter.Operator)
				}
			default:
				logging.Log.Warnf("Некорректный оператор фильтрации: %s для поля: %s", filter.Operator, filter.Field)
			}
		}
	}

	query = query.Limit(limit).Offset(offset)

	var users []models.Users
	if err := query.Find(&users).Error; err != nil {
		logging.Log.Errorf("Не удалось получить список пользователей %v", err)
		http.Error(w, fmt.Sprintf("Не удалось получить список пользователей: %v", err), 400)
		return
	}

	logging.Log.Debugf("Получены следующие пользователи: \n %v", users)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(users)

	logging.Log.Info("Запрос на получение списка пользователей успешно выполнен")

	return
}

func applyFilter(query *gorm.DB, field, value, operator string) *gorm.DB {
	switch operator {
	case "equals":
		return query.Where(fmt.Sprintf("%s = ?", field), value)
	case "contains":
		return query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
	case "startsWith":
		return query.Where(fmt.Sprintf("%s LIKE ?", field), value+"%")
	case "endsWith":
		return query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value)
	default:

		return query
	}
}
