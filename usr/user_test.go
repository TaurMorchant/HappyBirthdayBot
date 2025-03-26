package usr

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_SetBirthday(t *testing.T) {
	userDate, _ := time.Parse("02.01.2006", "13.06.2025")

	user := User{}
	user.SetBirthday(userDate, time.Now())

	assert.Equal(t, 13, user.birthday.GetDay())
	assert.Equal(t, "июня", user.birthday.GetMonth())
}

func Test_GetDaysBeforeBirthday_1(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "13.06.2025")
	userDate, _ := time.Parse("02.01.2006", "15.06.2025")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 2, user.GetDaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_2(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "15.06.2025")
	userDate, _ := time.Parse("02.01.2006", "15.07.2025")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 30, user.GetDaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_3(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "15.06.2025")
	userDate, _ := time.Parse("02.01.2006", "15.07.1111")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 30, user.GetDaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_4(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "31.12.2025")
	userDate, _ := time.Parse("02.01.2006", "01.01.2025")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 1, user.GetDaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_5(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "31.12.2025")
	userDate, _ := time.Parse("02.01.2006", "01.01.1111")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 1, user.GetDaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_6(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "11.11.2025")
	userDate, _ := time.Parse("02.01.2006", "11.11.1111")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 365, user.GetDaysBeforeBirthday())
}

func Test_FormattedString_1(t *testing.T) {
	userDate, _ := time.Parse("02.01.2006", "03.06.2025")

	user := User{Name: "Вася"}
	user.SetBirthday(userDate, time.Now())

	assert.Equal(t, "      Вася — 3  июня    ", user.FormattedString(10))
}

func Test_FormattedString_2(t *testing.T) {
	userDate, _ := time.Parse("02.01.2006", "13.09.2025")

	user := User{Name: "Вася"}
	user.SetBirthday(userDate, time.Now())

	assert.Equal(t, "      Вася — 13 сентября", user.FormattedString(10))
}
