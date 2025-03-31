package properties

//todo rename to properties?

import (
	"fmt"
	"github.com/magiconair/properties"
	"happy-birthday-bot/chat"
	res "happy-birthday-bot/resources"
)

var allowedUsers *properties.Properties
var allowedChats *properties.Properties

func init() {
	allowedUsersFile, err := res.ReadFile("allowedUsers.properties")
	if err != nil {
		panic(err)
	}
	allowedUsers = properties.MustLoadReader(allowedUsersFile, properties.UTF8)
	allowedChatsFile, err := res.ReadFile("allowedChats.properties")
	if err != nil {
		panic(err)
	}
	allowedChats = properties.MustLoadReader(allowedChatsFile, properties.UTF8)
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
