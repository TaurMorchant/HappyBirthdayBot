package handlers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/date"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
	"strings"
	"time"
	"unicode/utf8"
)

const Layout = "02.01.2006"

type JoinHandler struct {
}

func (h JoinHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	log.Printf("handle join command")
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	messageID := update.Message.MessageID

	users := sheets.Read()
	if user, ok := users.Get(usr.UserId(userID)); ok {
		msg := "Ты уже зарегистрирован! 😎\n\n"
		if len(user.Wishlist) == 0 {
			msg += "Хочешь задать свой вишлист? [/wishlist](/wishlist)"
		} else {
			msg += "Хочешь поменять свой вишлист? [/wishlist](/wishlist)"
		}

		bot.SendPic(chatID, msg, res.Cool)
		return nil
	}

	msg := "Отлично! Ответь на это сообщение вот так:\n\n`<Твое имя>, <Твоя дата рождения>`\n\n" +
		"Например:\n\n`Вася Пупкин, 25.03`\n\nили\n\n`Вася Пупкин, 25 марта`"
	bot.SendPicForceReply(chatID, msg, res.Waiting, messageID)

	WaitingForReplyHandlers.Add(userID, h)
	return nil
}

func (h JoinHandler) HandleReply(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()
	if _, ok := users.Get(usr.UserId(userID)); ok {
		bot.SendPic(chatID, "Ты уже зарегистрирован! 😎", res.Cool)
		return nil
	}

	name, birthDay, err := parseNameAndBirthday(update.Message.Text)
	if err != nil {
		return err
	}

	user := usr.User{Id: usr.UserId(userID), Name: name}
	user.SetBirthday2(birthDay, time.Now())

	users.Add(&user)
	sheets.Write(&users)

	bot.SendPic(chatID, "Поздравляю, теперь тебя отхеппибёздят! 🥳", res.Cool)

	return nil
}

func (h JoinHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}

//--------------------------------------------------------------------------------------------------

func parseNameAndBirthday(input string) (string, date.Birthday, error) {
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return "", date.Birthday{}, errors.New("Я тебя не понимаю. Ответь в виде\n\n`<Твое имя>, <Твоя дата рождения>`")
	}

	name := strings.TrimSpace(parts[0])
	if name == "" {
		return "", date.Birthday{}, errors.New("Ты забыл задать имя")
	}
	if utf8.RuneCountInString(name) > 10 {
		return "", date.Birthday{}, errors.New("Не наглей! Максимальная длина имени 10 символов!")
	}

	dateStr := strings.TrimSpace(parts[1])
	if dateStr == "" {
		return "", date.Birthday{}, errors.New("Ты забыл задать день рождения")
	}
	birthDay, err := parseDate(dateStr)
	return name, birthDay, err
}

func parseDate(input string) (date.Birthday, error) {
	birthdate, err := date.ParseBirthday(input)
	if err == nil {
		return birthdate, nil
	} else {
		dateTime, err := time.Parse(Layout, input)
		if err != nil {
			input = input + ".2000"
			dateTime, err = time.Parse(Layout, input)
			if err != nil {
				return date.Birthday{}, errors.New("Не понимаю твою дату. Вот корректные примеры:\n\n`Вася Пупкин, 25.03`\n\nили\n\n`Вася Пупкин, 25 марта`")
			}
		}
		return date.ToBirthday(dateTime), nil
	}
}
