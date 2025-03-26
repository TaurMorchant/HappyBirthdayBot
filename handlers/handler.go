package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/usr"
	"log"
	"time"
)

type IHandler interface {
	Handle(bot *bot.Bot, update tgbotapi.Update) error
	HandleReply(bot *bot.Bot, update tgbotapi.Update) error
}

var Handlers = map[string]IHandler{
	Test:      &TestHandler{},
	Start:     &StartHandler{},
	Join:      &JoinHandler{},
	Exit:      &ExitHandler{},
	List:      &ListHandler{},
	Reminders: &RemindHandler{},
}

var waitingHandlers = cache.New(30*time.Second, 1*time.Minute)

func WaitReply(userId usr.UserId, handler IHandler) {
	err := waitingHandlers.Add(fmt.Sprintf("%d", userId), handler, cache.DefaultExpiration)
	if err != nil {
		log.Panicln("Cannot add element to cache: ", err)
	}
}

func GetWaitingHandler(userId usr.UserId) IHandler {
	log.Println("All elements in cache: ", waitingHandlers.Items())
	cachedItem, ok := waitingHandlers.Get(fmt.Sprintf("%d", userId))
	if !ok {
		return nil
	}
	return cachedItem.(IHandler)
}

func RemoveWaitingHandler(userId usr.UserId) {
	waitingHandlers.Delete(fmt.Sprintf("%d", userId))
}
