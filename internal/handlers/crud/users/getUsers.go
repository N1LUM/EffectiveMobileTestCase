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

// GetUsers получает список пользователей с возможностью фильтрации и пагинации.
//
// Swagger: operationId=getUsers
// parameters:
//   - name: body
//     in: body
//     description: Параметры фильтрации пользователей.
//     required: false
//     schema:
//     "$ref": "#/definitions/UserGetListInput"
//
// responses:
//
//	'200':
//	  description: Успешное получение списка пользователей.
//	  schema:
//	    type: array
//	    items:
//	      "$ref": "#/definitions/User"
//	'400':
//	  description: Ошибка в запросе, например, неверный формат параметров фильтрации или ошибка получения данных.
//	  schema:
//	    type: object
//	    properties:
//	      error:
//	        type: string
//	        example: "Не удалось получить список пользователей: текст ошибки"
//
// definitions:
//
//	User:
//	  type: object
//	  properties:
//	    ID:
//	      type: string
//	      example: "550e8400-e29b-41d4-a716-446655440000"
//	    Name:
//	      type: string
//	      example: "Иван"
//	    Surname:
//	      type: string
//	      example: "Иванов"
//	    Patronymic:
//	      type: string
//	      example: "Иванович"
//	    Address:
//	      type: string
//	      example: "г. Москва, ул. Пушкина, д. Колотушкина"
//	    PassportSerie:
//	      type: string
//	      example: "1234"
//	    PassportNumber:
//	      type: string
//	      example: "567890"
//	    FullPassport:
//	      type: string
//	      example: "1234567890"
//	    CreatedAt:
//	      type: string
//	      format: date-time
//	      example: "2024-07-05T12:00:00Z"
//	    UpdatedAt:
//	      type: string
//	      format: date-time
//	      example: "2024-07-05T12:30:00Z"
//	UserGetListInput:
//	  type: object
//	  properties:
//	    Page:
//	      type: integer
//	      example: 1
//	    Limit:
//	      type: integer
//	      example: 10
//	    Filters:
//	      type: object
//	      properties:
//	        Filters:
//	          type: array
//	          items:
//	            type: object
//	            properties:
//	              Field:
//	                type: string
//	                example: "Name"
//	              Value:
//	                type: string
//	                example: "Иван"
//	              Operator:
//	                type: string
//	                example: "equals"
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
