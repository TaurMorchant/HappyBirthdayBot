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
	birthDay           date.Birthday
	daysBeforeBirthday int
	Wishlist           string
	Reminder30days     bool
	Reminder15days     bool
	BirthdayGreetings  bool
}

//---------------------------------------------------------------------------------------------------------------------

func (u *User) FormattedString(maxNameLength, maxMonthLength int) string {
	return fmt.Sprintf("%*s — %s", maxNameLength, u.Name, u.birthDay.ToStringWithFormatting(maxMonthLength))
}

func (u *User) SetBirthday(t time.Time, timeNow time.Time) {
	u.birthDay = date.ToBirthday(t)
	u.daysBeforeBirthday = u.calculateDaysBeforeBirthday(timeNow)
}

func (u *User) SetBirthday2(birthday date.Birthday, timeNow time.Time) {
	u.birthDay = birthday
	u.daysBeforeBirthday = u.calculateDaysBeforeBirthday(timeNow)
}

func (u *User) BirthDay() date.Birthday {
	return u.birthDay
}

func (u *User) DaysBeforeBirthday() int {
	return u.daysBeforeBirthday
}

//--------------------------------------------------------------

func (u *User) calculateDaysBeforeBirthday(timeNow time.Time) int {
	duration := timeNow.Sub(u.birthDay.CurrentYear())
	days := int(duration.Hours() / 24)

	if days >= 0 {
		return days
	} else {
		return 365 + days
	}
}
