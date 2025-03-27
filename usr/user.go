package usr

import (
	"fmt"
	"happy-birthday-bot/date"
	"time"
)

type UserId int64

type User struct {
	Id                 UserId
	Name               string
	birthday           date.Birthday
	daysBeforeBirthday int
	Wishlist           string
	Reminder30days     bool
	Reminder15days     bool
	BirthdayGreetings  bool
}

//---------------------------------------------------------------------------------------------------------------------

func (u *User) FormattedString(maxNameLength int) string {
	return fmt.Sprintf("%*s — %-11s", maxNameLength, u.Name, u.birthday.ToString())
}

func (u *User) SetBirthday(t time.Time, timeNow time.Time) {
	u.birthday = date.ToBirthday(t)
	u.daysBeforeBirthday = u.calculateDaysBeforeBirthday(timeNow)
}

func (u *User) Birthday() date.Birthday {
	return u.birthday
}

func (u *User) DaysBeforeBirthday() int {
	return u.daysBeforeBirthday
}

//--------------------------------------------------------------

func (u *User) calculateDaysBeforeBirthday(timeNow time.Time) int {
	timeNow = time.Date(u.birthday.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.UTC)

	duration := u.birthday.Sub(timeNow)
	days := int(duration.Hours() / 24)

	if days >= 0 {
		return days
	} else {
		return 365 + days
	}
}
