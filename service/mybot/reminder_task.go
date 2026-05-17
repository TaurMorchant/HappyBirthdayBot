package mybot

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"happy-birthday-bot/config"
	"happy-birthday-bot/db"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/usr"
	"log"
	"runtime/debug"
)

var mainChatId int64
var adminChatId int64

func StartReminderTask(bot *Bot) {
	mainChatId = config.GetInt64Property(config.MainChatIdProp)
	adminChatId = config.GetInt64Property(config.AdminChatIdProp)

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(config.GetStringProperty(config.ReminderTriggerCronProp), func() {
		isBirthdayComingUp(bot)
	})
	if err != nil {
		log.Panic("Не удалось запустить ReminderTask!", err)
	}

	c.Start()
}

//---------------------------------------------------------------------------------------

func isBirthdayComingUp(bot *Bot) {
	defer handlePanic(bot)
	log.Println("Check is birthday coming up")
	users := db.ReadUsers()
	for _, user := range users.AllUsers() {
		if user.DaysBeforeBirthday() == 0 && !user.BirthdayGreetings {
			handleBirthday(bot, user)
		} else if user.DaysBeforeBirthday() <= 14 && !user.Reminder15days {
			handle15Days(bot, user)
		} else if user.DaysBeforeBirthday() <= 30 && !user.Reminder30days {
			handle30Days(bot, user)
		} else if user.DaysBeforeBirthday() > 30 && user.BirthdayGreetings {
			user.BirthdayGreetings = false
			user.Reminder15days = false
			user.Reminder30days = false
		} else {
			continue
		}
		if err := db.UpdateFlags(user); err != nil {
			log.Panicf("Failed to persist flags for user %d: %v", user.Id, err)
		}
	}
}

func handleBirthday(bot *Bot, user *usr.User) {
	log.Printf("handleBirthday for user %v", user.Id)
	msg := fmt.Sprintf("Ура! Сегодня день рождения отмечает `%s`! 🎉", user.Name)
	bot.SendPic(mainChatId, msg, res.HappyBirthday)
	user.BirthdayGreetings = true
	user.Reminder15days = true
	user.Reminder30days = true
}

func handle15Days(bot *Bot, user *usr.User) {
	log.Printf("handle15Days for user %v", user.Id)
	birthdayChat := getBirthdayChat(user.Id)
	msg := fmt.Sprintf("Хочу напомнить, что и двух недель не осталось до момента, когда родится `%s`!", user.Name)
	if birthdayChat != nil {
		msg += fmt.Sprintf("\n\nЕсли ты всё ещё не присоединился к обсуждению подарка - самое время: %s", birthdayChat.ChatLink)
	}
	bot.SendPic(mainChatId, msg, res.Random)

	user.Reminder15days = true
	user.Reminder30days = true
}

func handle30Days(bot *Bot, user *usr.User) {
	log.Printf("handle30Days for user %v", user.Id)
	birthdayChat := getBirthdayChat(user.Id)
	msg := fmt.Sprintf("Псс, ребята! Уже меньше, чем через месяц, `%s` отмечает свой день рождения! Самое время подумать о подарке!", user.Name)
	if birthdayChat != nil {
		msg += fmt.Sprintf("\n\nЕсли хочешь обсудить, что подарим, залетай в чат: %s", birthdayChat.ChatLink)
	} else {
		msg += fmt.Sprintf("\n\nНо кажется @morchant ещё не завел чатик для обсуждения! Эй, пните его кто-нибудь!")
	}
	message := bot.SendPic(mainChatId, msg, res.Random)
	bot.PinMessage(mainChatId, message.MessageID)

	user.Reminder30days = true

	sendWishlistInBirthdayChat(bot, user)
}

func sendWishlistInBirthdayChat(bot *Bot, user *usr.User) {
	defer handlePanic(bot)

	msg := fmt.Sprintf("Это чат для обсуждения подарка. Именинник: `%v`. День рождения: %v.\n\n", user.Name, user.BirthDay().ToString())
	if len(user.Wishlist) == 0 {
		msg += fmt.Sprintf("Но похоже `%s` не составил виш лист :(", user.Name)
	} else {
		msg += fmt.Sprintf("Вишлист:\n\n```\n%s\n```", user.Wishlist)
	}

	birthdayChat := getBirthdayChat(user.Id)
	if birthdayChat != nil {
		message := bot.SendPic(birthdayChat.ChatId, msg, res.Wishlist)
		bot.PinMessage(birthdayChat.ChatId, message.MessageID)
	}
}

func getBirthdayChat(userId usr.UserId) *config.BirthdayChat {
	for _, bdChat := range config.BirthdayChats {
		if bdChat.UserId == int64(userId) {
			return &bdChat
		}
	}
	return nil
}

func handlePanic(bot *Bot) {
	if p := recover(); p != nil {
		log.Println("[PANIC] Panic was catch: ", p)
		log.Println(string(debug.Stack()))
		message := fmt.Sprintf("Поймана паника во время выполнения reminder task: %v", p)
		bot.SendText(adminChatId, message)
	}
}

//todo сделать команду пнутия, которая будет работать только через мою личку
