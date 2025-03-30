package date

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Birthday struct {
	day       int
	month     time.Month
	monthName string
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

func ParseBirthday(input string) (Birthday, error) {
	parts := strings.Fields(input)
	if len(parts) != 2 {
		return Birthday{}, errors.New("Invalid format")
	}

	dayStr := strings.TrimSpace(parts[0])
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return Birthday{}, err
	}

	monthStr := strings.TrimSpace(parts[1])
	month := getMonthNumber(monthStr)
	if month == -1 {
		return Birthday{}, errors.New("Invalid month")
	}

	dateTime, err := time.Parse("2006-1-2", fmt.Sprintf("2000-%d-%d", month, day))
	if err != nil {
		return Birthday{}, errors.New("Invalid date")
	}

	return ToBirthday(dateTime), nil
}

func ToBirthday(input time.Time) Birthday {
	day := input.Day()
	month := months[input.Month()]

	return Birthday{day: day, month: input.Month(), monthName: month}
}

func (b Birthday) CurrentYear() time.Time {
	return time.Date(time.Now().Year(), b.month, b.day, 0, 0, 0, 0, time.UTC)
}

func (b Birthday) ToString() string {
	return fmt.Sprintf("%d %s", b.day, b.monthName)
}

func (b Birthday) ToStringWithFormatting(maxMonthLength int) string {
	return fmt.Sprintf("%-2d %-*s", b.day, maxMonthLength, b.monthName)
}

func (b Birthday) MonthName() string {
	return b.monthName
}

func (b Birthday) Day() int {
	return b.day
}

func (b Birthday) Month() time.Month {
	return b.month
}

//-----------------------------------------------------------------------------------------------

func getMonthNumber(str string) time.Month {
	for i, v := range months {
		if v == str {
			return i
		}
	}
	return -1
}
