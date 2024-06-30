package db

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"test/internal/logging"
	"test/internal/models"
)

var PostgresClient *gorm.DB

func ConnectDB() {
	logging.Log.Info("Начало подключение к БД")

	err := godotenv.Load("./.env")
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Не удалось загрузить .env файл")
	}

	logging.Log.Debug("Загружен .env файл")

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	logging.Log.Debugf("Параметры подключения: username=%s, dbName=%s", username, dbName)

	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", username, password, dbName)

	PostgresClient, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Не удалось подключится к БД")
	}

	logging.Log.Debug("Открыто соединение с БД. Создана переменная PostgresClient для подключения к бд")

	logging.Log.Info("Успешное подключение к БД!")

	migrateModels()
}

func migrateModels() {
	logging.Log.Info("Начало миграции моделей в БД")
	// Миграция моделей
	err := PostgresClient.AutoMigrate(&models.Users{}, &models.Tasks{}, &models.UsersTasks{})
	if err != nil {
		logging.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Не удалось провести миграцию структуры БД")
	}
	logging.Log.Info("Успешная миграция!")
}
