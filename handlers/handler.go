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
	HandleCallback(bot *bot.Bot, update tgbotapi.Update) error
}

var Handlers = map[string]IHandler{
	Start:     &StartHandler{},
	Join:      &JoinHandler{},
	Exit:      &ExitHandler{},
	List:      &ListHandler{},
	Reminders: &RemindHandler{},
	Wishlist:  &WishlistHandler{},
}

// todo new package cache
var waitingForReplyHandlers = cache.New(30*time.Second, 1*time.Minute)

func WaitForReply(userId usr.UserId, handler IHandler) {
	err := waitingForReplyHandlers.Add(fmt.Sprintf("%d", userId), handler, cache.DefaultExpiration)
	if err != nil {
		log.Panicln("Cannot add element to reply cache: ", err)
	}
}

func GetWaitingForReplyHandler(userId usr.UserId) IHandler {
	log.Println("All elements in reply cache: ", waitingForReplyHandlers.Items())
	cachedItem, ok := waitingForReplyHandlers.Get(fmt.Sprintf("%d", userId))
	if !ok {
		return nil
	}
	return cachedItem.(IHandler)
}

func RemoveWaitingForReplyHandler(userId usr.UserId) {
	waitingForReplyHandlers.Delete(fmt.Sprintf("%d", userId))
}

//-----------------------------------------------------------------------------------------

type CallbackElement struct {
	UserId  int64
	Handler IHandler
}

// todo храть связку сообщения с пользователем, чтобы нельзя было нажать не на свою кнопку
var waitingForCallbackHandlers = cache.New(5*time.Minute, 10*time.Minute)

func WaitForCallback(messageId int, userId int64, handler IHandler) {
	callbackElement := CallbackElement{UserId: userId, Handler: handler}
	err := waitingForCallbackHandlers.Add(fmt.Sprintf("%d", messageId), callbackElement, cache.DefaultExpiration)
	if err != nil {
		log.Panicln("Cannot add element to callback cache: ", err)
	}
}

func GetWaitingForCallbackHandler(messageId int) *CallbackElement {
	log.Println("All elements in callback cache: ", waitingForCallbackHandlers.Items())
	cachedItem, ok := waitingForCallbackHandlers.Get(fmt.Sprintf("%d", messageId))
	if !ok {
		return nil
	}
	result := cachedItem.(CallbackElement)
	return &result
}

func RemoveWaitingForCallbackHandler(messageId int) {
	waitingForCallbackHandlers.Delete(fmt.Sprintf("%d", messageId))
}
