package logging

import "github.com/sirupsen/logrus"

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	// Настройка формата логирования
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		PrettyPrint:     true,
	})
	// Установка уровня логирования на DEBUG (включает INFO)
	Log.SetLevel(logrus.DebugLevel)

	Log.Info("Логгер инициализирован")
}
