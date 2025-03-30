package chat

import (
	res "happy-birthday-bot/resources"
	"log"
	"strconv"
	"strings"
)

type BirthdayChat struct {
	UserId   int64
	ChatLink string
	ChatId   int64
}

var BirthdayChats []BirthdayChat

func init() {
	rows, err := res.ReadCSV("birthdayChats.csv")
	if err != nil {
		log.Panic("Cannot read birthdayChats.csv", err)
	}
	for _, row := range rows {
		userIdStr := strings.TrimSpace(row[0])
		chatLink := strings.TrimSpace(row[1])
		chatIdStr := strings.TrimSpace(row[2])
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			log.Panic("Invalid userId", err)
		}
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			log.Panic("Invalid chatId", err)
		}
		birthdayChat := BirthdayChat{userId, chatLink, chatId}
		BirthdayChats = append(BirthdayChats, birthdayChat)
	}
}
