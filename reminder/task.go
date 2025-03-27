package reminder

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/robfig/cron/v3"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/handlers"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
)

var birthdayChats *properties.Properties

func init() {
	birthdayChats = properties.MustLoadFile("birthdayChats.properties", properties.UTF8)
}

func StartReminderTask(bot *bot.Bot) {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("* * 9 * * *", func() {
		isBirthdayComingUp(bot)
	})
	if err != nil {
		log.Panic("Не удалось запустить ReminderTask!", err)
	}

	c.Start()
}

func isBirthdayComingUp(bot *bot.Bot) {
	users := sheets.Read()
	isUpdateNeeded := false
	for _, user := range users.GetAllUsers() {
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

func handleBirthday(bot *bot.Bot, user *usr.User) {
	msg := fmt.Sprintf("Ура! Сегодня день рождения отмечает %s!", user.Name)
	bot.Send(handlers.MainChatId, msg)
	user.BirthdayGreetings = true
	user.Reminder15days = true
	user.Reminder30days = true
}

func handle15Days(bot *bot.Bot, user *usr.User) {
	chatLink := birthdayChats.GetString(fmt.Sprintf("%d", user.Id), "")
	msg := fmt.Sprintf("Хочу напомнить, что и двух недель не осталось до момента, когда родится `%s`!", user.Name)
	if chatLink != "" {
		msg += fmt.Sprintf("\n\nЕсли ты всё ещё не присоединился к обсуждению подарка - самое время: %s", chatLink)
	}
	bot.Send(handlers.MainChatId, msg)
	user.Reminder15days = true
	user.Reminder30days = true
}

func handle30Days(bot *bot.Bot, user *usr.User) {
	chatLink := birthdayChats.GetString(fmt.Sprintf("%d", user.Id), "")
	msg := fmt.Sprintf("Псс, ребята! Уже меньше, чем через месяц, `%s` отмечает свой день рождения! Самое время подумать о подарке!", user.Name)
	if chatLink != "" {
		msg += fmt.Sprintf("\n\nЕсли хочешь обсудить, что подарим, залетай в чат: %s", chatLink)
	} else {
		msg += fmt.Sprintf("\n\nНо кажется @morchant ещё не завел чатик для обсуждения! Ей, пните его кто-нибудь!")
	}
	bot.Send(handlers.MainChatId, msg)
	user.Reminder30days = true
}

//todo сделать команду пнутия, которая будет работать только через мою личку
