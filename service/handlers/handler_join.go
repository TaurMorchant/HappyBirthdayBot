package handlers

import (
	"errors"
	"fmt"
	"happy-birthday-bot/date"
	"happy-birthday-bot/db"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/usr"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const Layout = "02.01.2006"

const joinChangeDataButton = "join_change_data"
const joinChangeWishlistButton = "join_change_wishlist"
const joinAllOkButton = "join_all_ok"

type JoinHandler struct {
}

func (h JoinHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	log.Printf("handle my_birthday command")
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	messageID := update.Message.MessageID

	users := db.ReadUsers()
	if user, ok := users.Get(usr.UserId(userID)); ok {
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Поменять данные", joinChangeDataButton),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Задать вишлист", joinChangeWishlistButton),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Все ок", joinAllOkButton),
			),
		)
		msg1 := fmt.Sprintf("Я тебя уже знаю, `%v`! Твой день рождения `%v` 😎", user.Name, user.BirthDay().ToString())
		sentMessage := bot.SendPicWithKeyboard(chatID, msg1, res.Cool, &inlineKeyboard, messageID)
		WaitingForCallbackHandlers.Add(sentMessage.MessageID, CallbackElement{UserId: userID, Handler: h, OriginalMessageId: messageID})
		return nil
	}

	msg2 := "Отлично! Ответь на это сообщение вот так:\n\n`<Твое имя>, <Твоя дата рождения>`\n\n" +
		"Например:\n\n`Вася Пупкин, 25.03`\n\nили\n\n`Вася Пупкин, 25 марта`"
	bot.SendPicForceReply(chatID, msg2, res.Waiting, messageID)
	WaitingForReplyHandlers.Add(userID, h)
	return nil
}

func (h JoinHandler) HandleReply(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	name, birthDay, err := parseNameAndBirthday(update.Message.Text)
	if err != nil {
		return err
	}

	reactivated, err := db.ReactivateUser(usr.UserId(userID), name, birthDay)
	if err != nil {
		log.Panicf("Failed to reactivate user %d: %v", userID, err)
	}
	if reactivated {
		bot.SendPic(chatID, "С возвращением! Теперь тебя снова отхеппибёздят! 🥳", res.Cool)
		return nil
	}

	updated, err := db.UpdateUserData(usr.UserId(userID), name, birthDay)
	if err != nil {
		log.Panicf("Failed to update user %d: %v", userID, err)
	}
	if updated {
		bot.SendPic(chatID, "Данные обновлены! 👌", res.Cool)
		return nil
	}

	user := usr.User{Id: usr.UserId(userID), Name: name}
	user.SetBirthday2(birthDay, time.Now())
	if err := db.InsertUser(&user); err != nil {
		log.Panicf("Failed to insert user %d: %v", userID, err)
	}
	bot.SendPic(chatID, "Поздравляю, теперь тебя отхеппибёздят! 🥳", res.Cool)
	return nil
}

func (h JoinHandler) HandleCallback(bot *mybot.Bot, update tgbotapi.Update, callback CallbackElement) error {
	log.Println("Handle callback for JoinHandler")
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID

	switch update.CallbackQuery.Data {
	case joinChangeDataButton:
		msg := "Отлично! Ответь на это сообщение вот так:\n\n`<Твое имя>, <Твоя дата рождения>`\n\n" +
			"Например:\n\n`Вася Пупкин, 25.03`\n\nили\n\n`Вася Пупкин, 25 марта`"
		bot.SendPicForceReply(chatID, msg, res.Waiting, callback.OriginalMessageId)
		WaitingForReplyHandlers.Add(userID, h)
	case joinChangeWishlistButton:
		users := db.ReadUsers()
		if user, ok := users.Get(usr.UserId(userID)); ok {
			startWishlistFlow(bot, chatID, userID, callback.OriginalMessageId, user)
		}
	case joinAllOkButton:
		bot.SendPic(chatID, "Ну и славненько", res.Ok)
	default:
		bot.SendPic(chatID, "Ты откуда вообще взял эту кнопку, тут ее не должно быть!", res.Error)
	}
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
