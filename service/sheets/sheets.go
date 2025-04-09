package sheets

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"happy-birthday-bot/config"
	"happy-birthday-bot/date"
	"happy-birthday-bot/usr"
	"log"
	"strconv"
	"time"
)

var spreadsheetID string
var readRange string
var writeRange string

var srv *sheets.Service

func InitSpreadsheetService() {
	ctx := context.Background()

	spreadsheetConfig := config.SpreadsheetConfig()

	client := spreadsheetConfig.Client(ctx)

	var err error
	srv, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Panicf("Ошибка при создании сервиса: %v", err)
	}

	spreadsheetID = config.GetStringProperty(config.SpreadsheetIdProp)
	spreadsheetList := config.GetStringProperty(config.SpreadsheetListProp)
	readRange = spreadsheetList + "!A2:G30"
	writeRange = spreadsheetList + "!A2"
}

func Read() usr.Users {
	startTime := time.Now().Unix()
	var result usr.Users
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Panicf("Ошибка при чтении данных: %v", err)
	}

	for _, row := range resp.Values {
		user := readUser(row)
		if user != nil {
			result.Add(user)
		}
	}
	diff := time.Now().Unix() - startTime
	log.Println("Operation Read spent ", diff)
	return result
}

func Write(users *usr.Users) {
	startTime := time.Now().Unix()
	var values [][]interface{}
	for _, user := range users.AllUsers() {
		userRow := prepareUserRow(user)
		values = append(values, userRow)
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Записываем данные
	_, err := srv.Spreadsheets.Values.Clear(spreadsheetID, readRange, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		log.Panicf("Ошибка при удалении данных: %v", err)
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Panicf("Ошибка при записи данных: %v", err)
	}
	diff := time.Now().Unix() - startTime
	log.Println("Operation Write spent ", diff)
}

//-----------------------------

func readUser(row []interface{}) *usr.User {
	if len(row) == 0 {
		return nil
	}
	idStr := readRowValue(row, 0)
	name := readRowValue(row, 1)
	dateStr := readRowValue(row, 2)
	wishlist := readRowValue(row, 3)
	reminder30daysStr := readRowValue(row, 4)
	reminder15daysStr := readRowValue(row, 5)
	birthdayGreetingStr := readRowValue(row, 6)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Panicf("Ошибка при чтении ID. row: %s. Err: %s", row, err)
	}
	birthDay, err := date.ParseBirthday(dateStr)
	if err != nil {
		log.Panicf("Ошибка при чтении date. row: %s. Err: %s", row, err)
	}
	reminder30days := parseBool(reminder30daysStr)
	reminder15days := parseBool(reminder15daysStr)
	birthdayGreeting := parseBool(birthdayGreetingStr)

	result := &usr.User{Id: usr.UserId(id), Name: name, Wishlist: wishlist, Reminder30days: reminder30days, Reminder15days: reminder15days, BirthdayGreetings: birthdayGreeting}
	result.SetBirthday2(birthDay, time.Now())

	log.Println("read user: ", result)

	return result
}

func parseBool(str string) bool {
	result, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return result
}

func readRowValue(row []interface{}, i int) string {
	if i >= len(row) {
		return ""
	}
	return row[i].(string)
}

func prepareUserRow(user *usr.User) []interface{} {
	return []interface{}{user.Id, user.Name, user.BirthDay().ToString(), user.Wishlist, user.Reminder30days, user.Reminder15days, user.BirthdayGreetings}
}
