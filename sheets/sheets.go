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
const credentialsFile = "happybirthdaybot-454814-2dec5157295e.json"
const spreadsheetID = "1fb5ssf4Mp8HZ9aAFAOox9byQGUHstRub_5ssOdDoNro"
const readRange = "Data!A2:C30"
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
	var result usr.Users
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Ошибка при чтении данных: %v", err)
	}

	for _, row := range resp.Values {
		user := readUser(row)
		if user != nil {
			result.Add(*user)
		}
	}
	return result
}

func Write(users *usr.Users) {
	// Данные для записи (диапазон "A1:B1")
	var values [][]interface{}
	for _, user := range users.GetAllUsers() {
		userRow := []interface{}{user.Id, user.Name, date.FormatDate(user.GetBirthday().Time)}
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
}

//-----------------------------

func readUser(row []interface{}) *usr.User {
	if len(row) == 0 {
		return nil
	}
	idStr := row[0].(string)
	name := row[1].(string)
	dateStr := row[2].(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("Ошибка при чтении ID. row: %s. Err: %s", row, err)
	}
	dateTime, err := time.Parse(date.Layout, dateStr)
	if err != nil {
		log.Fatalf("Ошибка при чтении date. row: %s. Err: %s", row, err)
	}

	result := &usr.User{Id: usr.UserId(id), Name: name}
	result.SetBirthday(dateTime, time.Now())

	return result
}
