package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
)

type StartHandler struct {
}

func (h StartHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	msg := "Привет! Я — бот-отеппибёздывватель!\n" +
		"Я помогу вам организовывать групповые подарки на дни рождения!\n\n" +
		"Ты можешь использовать следующие комманды:\n\n"
	for _, command := range mybot.Commands {
		msg += fmt.Sprintf("[/%s] — %s\n", command.Command, command.Description)
	}
	msg += "\n\nА еще я пощу котиков `:3`"

	bot.SendPic(chatID, msg, res.Main)

	return nil
}

func (h StartHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error {
	return nil
}

func (h StartHandler) HandleCallback(*mybot.Bot, tgbotapi.Update) error {
	return nil
}
