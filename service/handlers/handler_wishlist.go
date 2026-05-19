package handlers

import (
	"fmt"
	"happy-birthday-bot/db"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/usr"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const viewOthersButton = "view_others"
const viewUserPrefix = "view_user_"

type WishlistHandler struct {
}

func (h WishlistHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	messageID := update.Message.MessageID

	users := db.ReadUsers()
	if user, ok := users.Get(usr.UserId(userID)); ok {
		startWishlistFlow(bot, chatID, userID, messageID, user)
	} else {
		bot.SendPic(chatID, "Кажется ты еще не зарегистрирован в программе! Зарегистрируйся при помощи команды [/my\\_birthday](/my_birthday)!", res.Suspicious)
	}

	return nil
}

func startWishlistFlow(bot *mybot.Bot, chatID int64, userID int64, messageID int, user *usr.User) {
	var msg string
	var firstButtonLabel string

	if len(user.Wishlist) == 0 {
		msg = "Похоже ты еще не составил свой вишлист! Самое время это сделать!"
		firstButtonLabel = "Хочу задать"
	} else {
		msg = fmt.Sprintf("У тебя сейчас такой вишлист:\n\n```\n%s\n```", user.Wishlist)
		firstButtonLabel = "Хочу поменять"
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(firstButtonLabel, okButton)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Хочу посмотреть другие", viewOthersButton)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Не, все норм", cancelButton)),
	)
	sentMessage := bot.SendPicWithKeyboard(chatID, msg, res.Wishlist, &inlineKeyboard, messageID)
	WaitingForCallbackHandlers.Add(sentMessage.MessageID, CallbackElement{UserId: userID, Handler: WishlistHandler{}, OriginalMessageId: messageID})
}

func (h WishlistHandler) HandleReply(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if err := db.UpdateWishlist(usr.UserId(userID), update.Message.Text); err != nil {
		log.Panicf("Failed to update wishlist for user %d: %v", userID, err)
	}
	bot.SendPic(chatID, "Вжух, вишлист обновлён! 👌", res.Vjuh)
	return nil
}

func (h WishlistHandler) HandleCallback(bot *mybot.Bot, update tgbotapi.Update, callback CallbackElement) error {
	log.Println("Handle callback for WishlistHandler")
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID
	data := update.CallbackQuery.Data

	if data == okButton {
		msg := "Напиши в ответ на это сообщение, что бы ты хотел получить в подарок?"
		messageToReplyId := callback.OriginalMessageId
		if messageToReplyId == 0 {
			messageToReplyId = messageID
		}
		bot.SendPicForceReply(chatID, msg, res.Wishlist, messageToReplyId)
		WaitingForReplyHandlers.Add(userID, h)
	} else if data == viewOthersButton {
		users := db.ReadUsers()
		var rows [][]tgbotapi.InlineKeyboardButton
		for _, u := range users.AllUsers() {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(u.Name, viewUserPrefix+strconv.FormatInt(int64(u.Id), 10)),
			))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
		sentMessage := bot.SendPicWithKeyboard(chatID, "Чей вишлист ты хочешь узнать?", res.Wishlist, &keyboard, 0)
		WaitingForCallbackHandlers.Add(sentMessage.MessageID, CallbackElement{UserId: userID, Handler: WishlistHandler{}})
	} else if strings.HasPrefix(data, viewUserPrefix) {
		targetUserIdStr := strings.TrimPrefix(data, viewUserPrefix)
		targetUserId, err := strconv.ParseInt(targetUserIdStr, 10, 64)
		if err != nil {
			bot.SendPic(chatID, "Что-то пошло не так 😕", res.Suspicious)
			return nil
		}
		users := db.ReadUsers()
		if targetUser, ok := users.Get(usr.UserId(targetUserId)); ok {
			var msg string
			if len(targetUser.Wishlist) == 0 {
				msg = fmt.Sprintf("`%s` не знает, чего хочет :(", targetUser.Name)
				bot.SendPic(chatID, msg, res.Sad)
			} else {
				msg = fmt.Sprintf("`%s` хочет:\n\n```\n%s\n```", targetUser.Name, targetUser.Wishlist)
				bot.SendPic(chatID, msg, res.Wishlist)
			}

		} else {
			bot.SendPic(chatID, "Такой участник не найден 🤔", res.Suspicious)
		}
	} else if data == cancelButton {
		bot.SendPic(chatID, "Океюшки", res.Ok)
	} else {
		bot.SendPic(chatID, "Ты откуда вообще взял эту кнопку, тут ее не должно быть!", res.Suspicious)
	}

	return nil
}
