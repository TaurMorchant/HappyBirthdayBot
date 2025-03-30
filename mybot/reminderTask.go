package mybot

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"happy-birthday-bot/chat"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
)

const MainChatId = 287959887
const EveryDay = "* * 9 * * *"
const Every10Sec = "*/10 * * * * *"

func StartReminderTask(bot *Bot) {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(Every10Sec, func() {
		isBirthdayComingUp(bot)
	})
	if err != nil {
		log.Panic("Не удалось запустить ReminderTask!", err)
	}

	c.Start()
}

func isBirthdayComingUp(bot *Bot) {
	users := sheets.Read()
	isUpdateNeeded := false
	for _, user := range users.AllUsers() {
		if user.DaysBeforeBirthday() == 0 && !user.BirthdayGreetings {
			handleBirthday(bot, user)
			isUpdateNeeded = true
			continue
		}
		if user.DaysBeforeBirthday() <= 14 && !user.Reminder15days {
			handle15Days(bot, user)
			isUpdateNeeded = true
			continue
		}
		if user.DaysBeforeBirthday() <= 30 && !user.Reminder30days {
			handle30Days(bot, user)
			isUpdateNeeded = true
			continue
		}
		if user.DaysBeforeBirthday() > 30 && user.BirthdayGreetings {
			//reset
			user.BirthdayGreetings = false
			user.Reminder15days = false
			user.Reminder30days = false
			isUpdateNeeded = true
		}

	}
	log.Println("isUpdateNeeded = ", isUpdateNeeded)
	if isUpdateNeeded {
		sheets.Write(&users)
	}
}

func handleBirthday(bot *Bot, user *usr.User) {
	msg := fmt.Sprintf("Ура! Сегодня день рождения отмечает '%s'!", user.Name)
	bot.SendPic(MainChatId, msg, res.HappyBirthday)
	user.BirthdayGreetings = true
	user.Reminder15days = true
	user.Reminder30days = true
}

func handle15Days(bot *Bot, user *usr.User) {
	birthdayChat := getBirthdayChat(user.Id)
	chatLink := birthdayChat.ChatLink
	msg := fmt.Sprintf("Хочу напомнить, что и двух недель не осталось до момента, когда родится `%s`!", user.Name)
	if chatLink != "" {
		msg += fmt.Sprintf("\n\nЕсли ты всё ещё не присоединился к обсуждению подарка - самое время: %s", chatLink)
	}
	message := bot.SendPic(MainChatId, msg, res.Random)
	bot.PinMessage(MainChatId, message.MessageID)

	user.Reminder15days = true
	user.Reminder30days = true
}

func handle30Days(bot *Bot, user *usr.User) {
	birthdayChat := getBirthdayChat(user.Id)
	msg := fmt.Sprintf("Псс, ребята! Уже меньше, чем через месяц, `%s` отмечает свой день рождения! Самое время подумать о подарке!", user.Name)
	if birthdayChat != nil {
		msg += fmt.Sprintf("\n\nЕсли хочешь обсудить, что подарим, залетай в чат: %s", birthdayChat.ChatLink)
	} else {
		msg += fmt.Sprintf("\n\nНо кажется @morchant ещё не завел чатик для обсуждения! Ей, пните его кто-нибудь!")
	}
	message := bot.SendPic(MainChatId, msg, res.Random)
	bot.PinMessage(MainChatId, message.MessageID)

	if len(user.Wishlist) == 0 {
		msg = fmt.Sprintf("Похоже `%s` не оставил виш лист :(", user.Name)
	} else {
		msg = fmt.Sprintf("Вот такой вишлист написал `%s`:\n\n```\n%s\n```", user.Name, user.Wishlist)
	}

	birthdayChat = getBirthdayChat(user.Id)
	if birthdayChat != nil {
		bot.SendText(birthdayChat.ChatId, msg)
	}

	user.Reminder30days = true
}

func getBirthdayChat(userId usr.UserId) *chat.BirthdayChat {
	for _, bdChat := range chat.BirthdayChats {
		if bdChat.UserId == int64(userId) {
			return &bdChat
		}
	}
	return nil
}

//todo сделать команду пнутия, которая будет работать только через мою личку
