package handlers

import (
	"happy-birthday-bot/cache"
	"happy-birthday-bot/mybot"
	"time"
)

var Handlers = map[string]IHandler{
	mybot.Start:     &StartHandler{},
	mybot.Join:      &JoinHandler{},
	mybot.Exit:      &ExitHandler{},
	mybot.List:      &ListHandler{},
	mybot.Reminders: &RemindHandler{},
	mybot.Wishlist:  &WishlistHandler{},
}

type CallbackElement struct {
	UserId            int64
	Handler           IHandler
	OriginalMessageId int
}

var WaitingForCallbackHandlers = cache.New[int, CallbackElement](5*time.Minute, 10*time.Minute)

var WaitingForReplyHandlers = cache.New[int64, IHandler](5*time.Minute, 10*time.Minute)
