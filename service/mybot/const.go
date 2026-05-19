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
	Select    = "select"
	Exec      = "exec"
	Remind    = "remind"
)

var Commands = []tgbotapi.BotCommand{
	{Command: List, Description: "Все дни рождения"},
	{Command: Reminders, Description: "Ближайшие дни рождения"},
	{Command: Join, Description: "Мой день рождения"},
	{Command: Wishlist, Description: "Настроить свой Wishlist"},
	{Command: Exit, Description: "Выйти из программы"},
}

var AdminCommands = []tgbotapi.BotCommand{
	{Command: Remind, Description: "[admin] Запустить задачу напоминаний вручную"},
	{Command: Exec, Description: "[admin] Выполнить SQL-запрос (INSERT/UPDATE/DELETE)"},
	{Command: Select, Description: "[admin] Выполнить SQL SELECT-запрос"},
	{Command: Answer, Description: "[admin] Отправить сообщение в основной чат"},
}
