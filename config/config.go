package config

import (
	"encoding/csv"
	"fmt"
	"github.com/magiconair/properties"
	"os"
	"strconv"
)

const MainChatIdProp = "mainChatId"
const AdminChatIdProp = "adminChatId"
const ReminderTriggerCronProp = "reminderTriggerCron"

var allowedUsers *properties.Properties
var allowedChats *properties.Properties
var applicationProperties *properties.Properties

func init() {
	allowedUsers = properties.MustLoadFile("./configs/allowedUsers.properties", properties.UTF8)
	allowedChats = properties.MustLoadFile("./configs/allowedChats.properties", properties.UTF8)
	applicationProperties = properties.MustLoadFile("./configs/application.properties", properties.UTF8)
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

//----------------------------------------------------------------------------------------------

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	return reader.ReadAll()
}
