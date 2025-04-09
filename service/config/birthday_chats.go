package config

import (
	"log"
	"strconv"
	"strings"
)

type BirthdayChat struct {
	UserId   int64
	Name     string
	ChatLink string
	ChatId   int64
}

var BirthdayChats []BirthdayChat

func initBirthdayChats(configsDir string) {
	rows, err := readCSV(configsDir + "/birthdayChats.csv")
	if err != nil {
		log.Panic("Cannot read birthdayChats.csv", err)
	}
	for _, row := range rows {
		userIdStr := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		chatLink := strings.TrimSpace(row[2])
		chatIdStr := strings.TrimSpace(row[3])
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			log.Panic("Invalid userId", err)
		}
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			log.Panic("Invalid chatId", err)
		}
		birthdayChat := BirthdayChat{userId, name, chatLink, chatId}
		BirthdayChats = append(BirthdayChats, birthdayChat)
	}
}
