package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/date"
	res "happy-birthday-bot/resources"
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
	messageID := update.Message.MessageID

	users := sheets.Read()
	if _, ok := users.Get(usr.UserId(userID)); ok {
		bot.SendPic(chatID, "Ты уже зарегистрирован!", res.Cool_cat)
		return nil
	}

	msg := "Отлично! Ответь на это сообщение вот так:\n\n`<Твое имя>, <дата рождения в формате DD.MM.YYYY>`\n\nНапример:\n\n`Вася Пупкин, 25.03.1990`"
	bot.SendTextForceReply(chatID, msg, messageID)

	WaitingForReplyHandlers.Add(userID, h)
	return nil
}

func (h JoinHandler) HandleReply(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()
	if _, ok := users.Get(usr.UserId(userID)); ok {
		bot.SendPic(chatID, "Ты уже зарегистрирован!", res.Cool_cat)
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

	bot.SendPic(chatID, "Поздравляю, теперь тебя отхеппибёздят!", res.Cool_cat)

	return nil
}

func (h JoinHandler) HandleCallback(*bot.Bot, tgbotapi.Update) error {
	return nil
}
