package date

import (
	"fmt"
	"time"
)

type Birthday struct {
	day       int
	monthName string
	time.Time
}

var months = map[time.Month]string{
	time.January:   "января",
	time.February:  "февраля",
	time.March:     "марта",
	time.April:     "апреля",
	time.May:       "мая",
	time.June:      "июня",
	time.July:      "июля",
	time.August:    "августа",
	time.September: "сентября",
	time.October:   "октября",
	time.November:  "ноября",
	time.December:  "декабря",
}

//-----------------------------------------------------

func ToBirthday(input time.Time) Birthday {
	currentYear := time.Now().Year()

	date := time.Date(currentYear, input.Month(), input.Day(), 0, 0, 0, 0, time.UTC)
	day := input.Day()
	month := months[input.Month()]

	return Birthday{day: day, monthName: month, Time: date}
}

func (b Birthday) ToString(maxMonthLength int) string {
	return fmt.Sprintf("%-2d %-*s", b.day, maxMonthLength, b.monthName)
}

func (b Birthday) MonthName() string {
	return b.monthName
}

func (b Birthday) GetDay() int {
	return b.day
}
