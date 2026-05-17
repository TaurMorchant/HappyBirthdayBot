package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/db"
	"happy-birthday-bot/mybot"
	"log"
	"time"
)

type DbTestHandler struct{}

func (h DbTestHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	log.Println("Handle dbtest command")
	chatID := update.Message.Chat.ID

	text := fmt.Sprintf("Запись из Telegram, %s", time.Now().Format("2006-01-02 15:04:05"))
	if err := db.InsertNote(text); err != nil {
		return err
	}

	notes, err := db.GetAllNotes()
	if err != nil {
		return err
	}

	msg := "SQLite тест — последние записи:\n```\n"
	for _, n := range notes {
		msg += fmt.Sprintf("#%d | %s | %s\n", n.ID, n.CreatedAt, n.Text)
	}
	msg += "```"

	bot.SendText(chatID, msg)
	return nil
}

func (h DbTestHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error { return nil }
func (h DbTestHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
