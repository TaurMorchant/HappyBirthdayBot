package impl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"happy-birthday-bot/date"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
	"strconv"
	"time"
)

type JoinHandler struct {
}

var joinRequests = cache.New(30*time.Second, 30*time.Second)

func (h JoinHandler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("handle join command")
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	msg := "Отлично! Ответь на это сообщение вот так:\n\n`<Твое имя>, <дата рождения в формате DD.MM.YYYY>`\n\nНапример:\n\n`Вася Пупкин, 25.03.1990`"
	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = tgbotapi.ModeMarkdown
	message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true} // Принудительный reply-режим

	sentMess, err := bot.Send(message)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
	log.Printf("ID отправленного сообщения: %s", sentMess.MessageID)
	joinRequests.Add(strconv.FormatInt(userID, 10), "wait for reply", cache.DefaultExpiration)
	log.Println("join requests: ", joinRequests)
}

func HandleReply(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	log.Printf("handle reply command")
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	_, ok := joinRequests.Get(strconv.FormatInt(userID, 10))
	if ok {
		users := sheets.Read()
		if _, ok := users.Get(usr.UserId(userID)); ok {
			bot.Send(tgbotapi.NewMessage(chatID, "Ты уже зарегистрирован!"))
			return nil
		}

		name, birthdate, err := date.ParseNameAndBirthdate(update.Message.Text)
		if err != nil {
			//log.Printf(err.Error())
			//bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
			return err
		}

		users.Add(usr.User{Id: usr.UserId(userID), Name: name, Birthday: birthdate})
		sheets.Write(&users)

		bot.Send(tgbotapi.NewMessage(chatID, "Поздравляю, теперь тебя отхеппибёздят!"))
		joinRequests.Delete(strconv.FormatInt(userID, 10))
	}
	return nil
}
