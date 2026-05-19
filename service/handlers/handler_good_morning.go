package handlers

import (
	"happy-birthday-bot/mybot"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const catAPIURL = "https://genrandom.com/api/cat"

type GoodMorningHandler struct {
}

func (h GoodMorningHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	resp, err := http.Get(catAPIURL)
	if err != nil {
		log.Printf("[ERROR] good_morning: fetch cat failed: %v", err)
		bot.SendText(chatID, "С добрым утром! 🐱")
		return nil
	}
	defer resp.Body.Close()

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] good_morning: read cat body failed: %v", err)
		bot.SendText(chatID, "С добрым утром! 🐱")
		return nil
	}

	bot.SendPicBytes(chatID, "С добрым утром!", imgBytes)
	return nil
}

func (h GoodMorningHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error {
	return nil
}

func (h GoodMorningHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
