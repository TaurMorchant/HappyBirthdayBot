package restrictions

//todo rename to properties?

import (
	"fmt"
	"github.com/magiconair/properties"
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
	return allowedChats.GetString(fmt.Sprintf("%d", chatId), "") != ""
}
