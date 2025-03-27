package handlers

import (
	"happy-birthday-bot/cache"
	"time"
)

var Handlers = map[string]IHandler{
	Start:     &StartHandler{},
	Join:      &JoinHandler{},
	Exit:      &ExitHandler{},
	List:      &ListHandler{},
	Reminders: &RemindHandler{},
	Wishlist:  &WishlistHandler{},
}

type CallbackElement struct {
	UserId  int64
	Handler IHandler
}

var WaitingForCallbackHandlers = cache.New[int, CallbackElement](30*time.Second, 1*time.Minute)

var WaitingForReplyHandlers = cache.New[int64, IHandler](30*time.Second, 1*time.Minute)
