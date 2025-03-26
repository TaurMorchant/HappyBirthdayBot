package reminder

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"log"
	"time"
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
	timeNow := time.Now()
	users := sheets.Read()
	for _, user := range users.GetAllUsers() {
		if user.GetBirthday().Before(timeNow) {
			continue
		}
		days := user.GetDaysBeforeBirthday()
		log.Println("days = ", days)
		if days < 30 {
			message := tgbotapi.NewMessage(287959887, fmt.Sprintf("До дня рождения %s осталось меньше 30 дней!"))
			bot.SendWithEH(message)
		}
		break
	}
}
