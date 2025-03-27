package sheets

import (
	"context"
	"golang.org/x/oauth2/google"
	"happy-birthday-bot/date"
	"happy-birthday-bot/usr"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Укажи путь к JSON-ключу
const credentialsFile = "C:\\go_modules\\happy_birthday_bot\\happybirthdaybot-454814-2dec5157295e.json"
const spreadsheetID = "1fb5ssf4Mp8HZ9aAFAOox9byQGUHstRub_5ssOdDoNro"
const readRange = "Data!A2:G30"
const writeRange = "Data!A2"

var srv *sheets.Service

func init() {
	// Создаём контекст
	ctx := context.Background()

	// Загружаем учетные данные из JSON-файла
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Не удалось прочитать файл ключа: %v", err)
	}

	// Настраиваем клиента
	config, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Ошибка при настройке JWT: %v", err)
	}

	client := config.Client(ctx)

	// Подключаемся к API
	srv, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Ошибка при создании сервиса: %v", err)
	}
}

func Read() usr.Users {
	startTime := time.Now().Unix()
	var result usr.Users
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Ошибка при чтении данных: %v", err)
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
	for _, user := range users.GetAllUsers() {
		userRow := prepareUserRow(user)
		values = append(values, userRow)
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Записываем данные
	_, err := srv.Spreadsheets.Values.Clear(spreadsheetID, readRange, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		log.Fatalf("Ошибка при удалении данных: %v", err)
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Ошибка при записи данных: %v", err)
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
		log.Fatalf("Ошибка при чтении ID. row: %s. Err: %s", row, err)
	}
	dateTime, err := time.Parse(date.Layout, dateStr)
	if err != nil {
		log.Fatalf("Ошибка при чтении date. row: %s. Err: %s", row, err)
	}
	reminder30days := parseBool(reminder30daysStr)
	reminder15days := parseBool(reminder15daysStr)
	birthdayGreeting := parseBool(birthdayGreetingStr)

	result := &usr.User{Id: usr.UserId(id), Name: name, Wishlist: wishlist, Reminder30days: reminder30days, Reminder15days: reminder15days, BirthdayGreetings: birthdayGreeting}
	result.SetBirthday(dateTime, time.Now())

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
	return []interface{}{user.Id, user.Name, date.FormatDate(user.Birthday().Time), user.Wishlist, user.Reminder30days, user.Reminder15days, user.BirthdayGreetings}
}
