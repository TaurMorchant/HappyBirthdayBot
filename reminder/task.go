package reminder

import (
	"github.com/robfig/cron/v3"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"log"
)

func StartReminderTask(bot *bot.Bot) {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("*/10 * * * * *", func() { // Каждые 10 секунд
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
	for _, user := range *users.GetAllUsers() {
		if user.DaysBeforeBirthday() == 365 && !user.BirthdayGreetings {
			//поздравить с др
			log.Println("Поздравить с ДР ", user.Name)
			user.BirthdayGreetings = true
			isUpdateNeeded = true
		}
		if user.DaysBeforeBirthday() <= 15 && !user.Reminder15days {
			//напомнить о 15 днях
			log.Printf("Напоминаю, до ДР %s осталось %d дней", user.Name, user.DaysBeforeBirthday())
			user.Reminder15days = true
			isUpdateNeeded = true
		}
		if user.DaysBeforeBirthday() <= 30 && !user.Reminder30days {
			//напомнить о 30 днях
			log.Printf("Эй, до ДР %s осталось %d дней", user.Name, user.DaysBeforeBirthday())
			user.Reminder30days = true
			isUpdateNeeded = true
		}
	}
	log.Println("isUpdateNeeded = ", isUpdateNeeded)
	if isUpdateNeeded {
		sheets.Write(&users)
	}
}
