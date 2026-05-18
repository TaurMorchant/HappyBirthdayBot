package mybot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	Start     = "start"
	List      = "list"
	Reminders = "reminders"
	Join      = "my_birthday"
	Wishlist  = "wishlist"
	Exit      = "exit"
	Answer    = "a"
)

var Commands = []tgbotapi.BotCommand{
	{Command: List, Description: "Все дни рождения"},
	{Command: Reminders, Description: "Ближайшие дни рождения"},
	{Command: Join, Description: "Мой день рождения"},
	{Command: Wishlist, Description: "Настроить свой Wishlist"},
	{Command: Exit, Description: "Выйти из программы"},
}
