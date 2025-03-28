package date

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

const Layout = "02.01.2006"

func ParseNameAndBirthdate(input string) (name string, birthdate time.Time, err error) {
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return "", time.Time{}, errors.New("Я тебя не понимаю. Ответь в виде\n\n`<Твое имя>, <дата рождения в формате DD.MM.YYYY>`")
	}

	name = strings.TrimSpace(parts[0])
	if name == "" {
		return "", time.Time{}, errors.New("Ты забыл задать имя")
	}
	if utf8.RuneCountInString(name) > 10 {
		return "", time.Time{}, errors.New("Не наглей! Максимальная длина имени 10 символов!")
	}

	dateStr := strings.TrimSpace(parts[1])
	birthdate, err = time.Parse(Layout, dateStr)
	if err != nil {
		return "", time.Time{}, errors.New("Не понимаю твою дату. Напиши ее пожалуйста в виде `DD.MM.YYYY`")
	}

	return name, birthdate, nil
}

func FormatDate(t time.Time) string {
	return t.Format(Layout)
}
