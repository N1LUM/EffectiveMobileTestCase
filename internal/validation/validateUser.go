package validation

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"test/internal/db"
	"test/internal/logging"
	"test/internal/models"
	"unicode"
	"unicode/utf8"
)

func ValidateUser(user *models.Users) error {
	logging.Log.Info("Начало валидации данных пользователя")

	logging.Log.WithFields(logrus.Fields{
		"user_id":             user.ID,
		"user_name":           user.Name,
		"user_surname":        user.Surname,
		"user_patronymic":     user.Patronymic,
		"user_address":        user.Address,
		"user_passportSerie":  user.PassportSerie,
		"user_pussportNumber": user.PassportNumber,
		"user_fullPassport":   user.FullPassport,
		"user_createdAt":      user.CreatedAt,
		"user_updatedAt":      user.UpdatedAt,
	}).Debug("В валидацию пришли следующие данные")

	if err := validateUserName(user); err != nil {
		logging.Log.Error("Валидация имени пользователя провалилась")
		return err
	}
	if err := validateUserSurname(user); err != nil {
		logging.Log.Error("Валидация фамилии пользователя провалилась")
		return err
	}
	if err := validateUserPatronymic(user); err != nil {
		logging.Log.Error("Валидация отчества пользователя провалилась")
		return err
	}
	if err := validateAddress(user); err != nil {
		logging.Log.Error("Валидация адреса пользователя провалилась")
		return err
	}
	if err := validatePassportSerie(user); err != nil {
		logging.Log.Error("Валидация серии паспорта пользователя провалилась")
		return err
	}
	if err := validatePassportNumber(user); err != nil {
		logging.Log.Error("Валидация номера паспорта пользователя провалилась")
		return err
	}
	if err := validateFullPassport(user); err != nil {
		logging.Log.Error("Валидация полного номера паспорта пользователя провалилась")
		return err
	}

	logging.Log.Info("Валидация пользователя успешно завершена!")

	return nil
}

func validateUserName(user *models.Users) error {
	logging.Log.Debugf("Длина имени=%d", utf8.RuneCountInString(user.Name))

	if utf8.RuneCountInString(user.Name) == 0 {
		return fmt.Errorf("У пользователя отсутствует имя!")
	}
	if utf8.RuneCountInString(user.Name) < 2 {
		return fmt.Errorf("Длинна имени пользователя должна быть не меньше 2х символов! Имя=%s", user.Name)
	}
	if !onlyLetters(user.Name) {
		return fmt.Errorf("Имя пользователя должно содержать только буквы! Имя=%s", user.Name)
	}

	user.Name = normalizedString(user.Name)

	return nil
}
func validateUserSurname(user *models.Users) error {
	logging.Log.Debugf("Длина фамилии=%d", utf8.RuneCountInString(user.Surname))

	if utf8.RuneCountInString(user.Surname) == 0 {
		return fmt.Errorf("У пользователя отсутствует фамилия!")
	}
	if utf8.RuneCountInString(user.Surname) < 2 {
		return fmt.Errorf("Длина фамилии пользователя должна быть не меньше 2-х символов! Фамилия=%s", user.Surname)
	}
	if !onlyLetters(user.Surname) {
		return fmt.Errorf("Фамилия пользователя должна содержать только буквы! Фамилия=%s", user.Surname)
	}

	user.Surname = normalizedString(user.Surname)

	return nil
}
func validateUserPatronymic(user *models.Users) error {
	logging.Log.Debugf("Длина отчества=%d", utf8.RuneCountInString(user.Patronymic))

	if utf8.RuneCountInString(user.Patronymic) > 21 {
		return fmt.Errorf("Слишком длинное отчество пользователя! Такого не существует! Отчество=%s", user.Patronymic)
	}
	if !onlyLetters(user.Patronymic) {
		return fmt.Errorf("Отчество пользователя должно содержать только буквы! Отчество=%s", user.Patronymic)
	}

	user.Patronymic = normalizedString(user.Patronymic)

	return nil
}
func validateAddress(user *models.Users) error {
	for _, r := range user.Address {
		// Если символ не является буквой, цифрой или одним из разрешённых специальных символов, возвращаем false
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune(" .,-/", r) {
			return fmt.Errorf("Адрес содержит запрещенные символы! Адрес=%s", user.Address)
		}
	}
	return nil
}
func validatePassportSerie(user *models.Users) error {
	logging.Log.Debugf("Длина серии паспорта=%d", utf8.RuneCountInString(user.PassportSerie))

	if utf8.RuneCountInString(user.PassportSerie) != 4 {
		return fmt.Errorf("Длина серии паспорта должна ровняться 4! Серия паспорта=%s", user.PassportSerie)
	}
	if !onlyNumbers(user.PassportSerie) {
		return fmt.Errorf("Серия паспорта должна содержать только цифры! Серия паспорта=%s", user.PassportSerie)
	}
	return nil
}
func validatePassportNumber(user *models.Users) error {
	logging.Log.Debugf("Длина номера паспорта=%d", utf8.RuneCountInString(user.PassportNumber))

	if utf8.RuneCountInString(user.PassportNumber) != 6 {
		return fmt.Errorf("Длина номера паспорта должна ровняться 6! Номер паспорта=%s", user.PassportNumber)
	}
	if !onlyNumbers(user.PassportNumber) {
		return fmt.Errorf("Номер паспорта должен содержать только цифры! Номер паспорта=%s", user.PassportNumber)
	}

	return nil
}

func validateFullPassport(user *models.Users) error {
	var count int64
	err := db.PostgresClient.Model(&models.Users{}).Where("full_passport = ?", user.FullPassport).Count(&count).Error
	if err != nil {
		return fmt.Errorf("Ошибка при выполнении запроса: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("Пользователь с таким паспортом уже существует: %s", user.FullPassport)
	}

	return nil
}

// Делаем первую букву заглавной
func normalizedString(str string) string {
	str = strings.ToLower(str)
	strRune := []rune(str)
	strRune[0] = unicode.ToUpper(strRune[0])
	return string(strRune)
}

// Проверяем, что все символы, только буквы
func onlyLetters(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// Проверяем, что все символы, только цифры
func onlyNumbers(str string) bool {
	for _, r := range str {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}
