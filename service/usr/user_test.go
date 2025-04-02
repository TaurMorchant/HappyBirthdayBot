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

	assert.Equal(t, 13, user.birthDay.Day())
	assert.Equal(t, "июня", user.birthDay.MonthName())
}

func Test_GetDaysBeforeBirthday_1(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "13.06.2025")
	userDate, _ := time.Parse("02.01.2006", "15.06.2025")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 2, user.DaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_2(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "15.06.2025")
	userDate, _ := time.Parse("02.01.2006", "15.07.2025")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 30, user.DaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_3(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "15.06.2025")
	userDate, _ := time.Parse("02.01.2006", "15.07.1111")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 30, user.DaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_4(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "31.12.2025")
	userDate, _ := time.Parse("02.01.2006", "01.01.2025")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 1, user.DaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_5(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "31.12.2025")
	userDate, _ := time.Parse("02.01.2006", "01.01.1111")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 1, user.DaysBeforeBirthday())
}

func Test_GetDaysBeforeBirthday_6(t *testing.T) {
	currentDate, _ := time.Parse("02.01.2006", "11.11.2025")
	userDate, _ := time.Parse("02.01.2006", "11.11.1111")

	user := User{}
	user.SetBirthday(userDate, currentDate)

	assert.Equal(t, 0, user.DaysBeforeBirthday())
}

func Test_FormattedString_1(t *testing.T) {
	userDate, _ := time.Parse("02.01.2006", "03.06.2025")

	user := User{Name: "Вася"}
	user.SetBirthday(userDate, time.Now())

	assert.Equal(t, "      Вася — 3  июня    ", user.FormattedString(10, 8))
}

func Test_FormattedString_2(t *testing.T) {
	userDate, _ := time.Parse("02.01.2006", "13.09.2025")

	user := User{Name: "Вася"}
	user.SetBirthday(userDate, time.Now())

	assert.Equal(t, "      Вася — 13 сентября", user.FormattedString(10, 8))
}

func Test_calculateDaysBetween(t *testing.T) {
	type args struct {
		time1 time.Time
		time2 time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "1", args: args{time1: time.Date(1999, time.April, 3, 0, 0, 0, 0, time.Local), time2: time.Date(1999, time.April, 3, 0, 0, 0, 0, time.Local)}, want: 0},
		{name: "2", args: args{time1: time.Date(1999, time.April, 3, 0, 0, 0, 0, time.Local), time2: time.Date(1999, time.April, 4, 0, 0, 0, 0, time.Local)}, want: 1},
		{name: "3", args: args{time1: time.Date(1999, time.April, 3, 0, 0, 0, 0, time.Local), time2: time.Date(1999, time.May, 3, 0, 0, 0, 0, time.Local)}, want: 30},
		{name: "4", args: args{time1: time.Date(1999, time.April, 3, 0, 0, 0, 0, time.Local), time2: time.Date(2000, time.April, 2, 0, 0, 0, 0, time.Local)}, want: 365},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, calculateDaysBetween(tt.args.time1, tt.args.time2), "calculateDaysBetween(%v, %v)", tt.args.time1, tt.args.time2)
		})
	}
}
