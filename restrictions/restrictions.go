package restrictions

//todo rename to properties?

import (
	"fmt"
	"github.com/magiconair/properties"
	"happy-birthday-bot/chat"
)

var allowedUsers *properties.Properties
var allowedChats *properties.Properties

func init() {
	allowedUsers = properties.MustLoadFile("allowedUsers.properties", properties.UTF8)
	allowedChats = properties.MustLoadFile("allowedChats.properties", properties.UTF8)
}

func IsUserAllowed(userId int64) bool {
	return allowedUsers.GetString(fmt.Sprintf("%d", userId), "") != ""
}

func IsChatAllowed(chatId int64) bool {
	for _, chat := range chat.BirthdayChats {
		if chat.ChatId == chatId {
			return true
		}
	}
	return allowedChats.GetString(fmt.Sprintf("%d", chatId), "") != ""
}
