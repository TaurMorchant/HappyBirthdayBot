package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/date"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
	"time"
)

type JoinHandler struct {
}

func (h JoinHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	log.Printf("handle join command")
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()
	if _, ok := users.Get(usr.UserId(userID)); ok {
		bot.SendWithEH(tgbotapi.NewMessage(chatID, "Ты уже зарегистрирован!"))
		return nil
	}

	msg := "Отлично! Ответь на это сообщение вот так:\n\n`<Твое имя>, <дата рождения в формате DD.MM.YYYY>`\n\nНапример:\n\n`Вася Пупкин, 25.03.1990`"
	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = tgbotapi.ModeMarkdown
	message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true} // Принудительный reply-режим

	bot.SendWithEH(message)

	WaitForReply(usr.UserId(userID), h)
	return nil
}

func (h JoinHandler) HandleReply(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()
	if _, ok := users.Get(usr.UserId(userID)); ok {
		bot.SendWithEH(tgbotapi.NewMessage(chatID, "Ты уже зарегистрирован!"))
		return nil
	}

	name, birthdate, err := date.ParseNameAndBirthdate(update.Message.Text)
	if err != nil {
		return err
	}

	user := usr.User{Id: usr.UserId(userID), Name: name}
	user.SetBirthday(birthdate, time.Now())

	users.Add(&user)
	sheets.Write(&users)

	bot.SendWithEH(tgbotapi.NewMessage(chatID, "Поздравляю, теперь тебя отхеппибёздят!"))

	return nil
}
