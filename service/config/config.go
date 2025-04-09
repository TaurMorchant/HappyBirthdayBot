package config

import (
	"encoding/csv"
	"fmt"
	"github.com/magiconair/properties"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"strconv"
)

const CONFIG_PATH_TMP = "../configs-test/"

const MainChatIdProp = "mainChatId"
const AdminChatIdProp = "adminChatId"
const ReminderTriggerCronProp = "reminderTriggerCron"
const SpreadsheetListProp = "spreadsheetList"
const SpreadsheetIdProp = "spreadsheetID"

var allowedUsers *properties.Properties
var allowedChats *properties.Properties
var applicationProperties *properties.Properties
var spreadsheetConfig *jwt.Config

func InitConfigs() {
	allowedUsers = properties.MustLoadFile(CONFIG_PATH_TMP+"allowedUsers.properties", properties.UTF8)
	allowedChats = properties.MustLoadFile(CONFIG_PATH_TMP+"allowedChats.properties", properties.UTF8)
	applicationProperties = properties.MustLoadFile(CONFIG_PATH_TMP+"application.properties", properties.UTF8)

	// Загружаем учетные данные из JSON-файла
	data, err := os.ReadFile(CONFIG_PATH_TMP + "happybirthdaybot-454814-2dec5157295e.json")
	if err != nil {
		log.Panicf("Не удалось прочитать файл ключа: %v", err)
	}

	// Настраиваем клиента
	spreadsheetConfig, err = google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	if err != nil {
		log.Panicf("Ошибка при настройке JWT: %v", err)
	}

	initBirthdayChats()
}

func IsUserAllowed(userId int64) bool {
	return allowedUsers.GetString(fmt.Sprintf("%d", userId), "") != ""
}

func IsChatAllowed(chatId int64) bool {
	for _, chat := range BirthdayChats {
		if chat.ChatId == chatId {
			return true
		}
	}
	return allowedChats.GetString(fmt.Sprintf("%d", chatId), "") != ""
}

func GetStringProperty(key string) string {
	result, ok := applicationProperties.Get(key)
	if !ok {
		msg := fmt.Sprintf("No required property found for key: %s", key)
		panic(msg)
	}
	return result
}

func GetInt64Property(key string) int64 {
	value, ok := applicationProperties.Get(key)
	if !ok {
		msg := fmt.Sprintf("No required property found for key: %s", key)
		panic(msg)
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid property value: %s. int64 required", key)
		panic(msg)
	}
	return result
}

func SpreadsheetConfig() *jwt.Config {
	return spreadsheetConfig
}

//----------------------------------------------------------------------------------------------

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Cannot close file")
		}
	}(file)

	reader := csv.NewReader(file)

	return reader.ReadAll()
}
