package usr

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"happy-birthday-bot/date"
	"testing"
	"time"
)

func Test_GetMaxNameLength(t *testing.T) {
	users := Users{}
	users.Add(&User{Name: "Qwert"})
	users.Add(&User{Name: "Qwerty"})
	users.Add(&User{Name: "Qwer"})

	assert.Equal(t, 6, users.GetMaxNameLength())
}

func Test_sort(t *testing.T) {
	users := Users{}
	user1 := User{Name: "1", birthDay: date.ToBirthday(time.Date(1991, 11, 4, 0, 0, 0, 0, time.UTC))}
	user2 := User{Name: "2", birthDay: date.ToBirthday(time.Date(1991, 12, 5, 0, 0, 0, 0, time.UTC))}
	user3 := User{Name: "3", birthDay: date.ToBirthday(time.Date(1991, 5, 7, 0, 0, 0, 0, time.UTC))}
	users.Add(&user1)
	users.Add(&user2)
	users.Add(&user3)

	users.sort()

	expected := []*User{&user3, &user1, &user2}

	assert.True(t, slices.Equal(expected, users.AllUsers()))
}

func Test_sortByDaysBeforeBirthday(t *testing.T) {
	currentDate := time.Date(time.Now().Year(), 5, 1, 0, 0, 0, 0, time.UTC)

	users := Users{}
	user1 := User{Name: "1"}
	user1.SetBirthday(time.Date(1991, 1, 4, 0, 0, 0, 0, time.UTC), currentDate)
	user2 := User{Name: "2"}
	user2.SetBirthday(time.Date(1991, 2, 13, 0, 0, 0, 0, time.UTC), currentDate)
	user3 := User{Name: "3"}
	user3.SetBirthday(time.Date(1991, 5, 7, 0, 0, 0, 0, time.UTC), currentDate)
	user4 := User{Name: "4"}
	user4.SetBirthday(time.Date(1991, 7, 11, 0, 0, 0, 0, time.UTC), currentDate)
	user5 := User{Name: "5"}
	user5.SetBirthday(time.Date(1991, 12, 2, 0, 0, 0, 0, time.UTC), currentDate)
	users.Add(&user1)
	users.Add(&user2)
	users.Add(&user3)
	users.Add(&user4)
	users.Add(&user5)

	result := users.sortByDaysBeforeBirthday()

	expected := []*User{&user3, &user4, &user5, &user1, &user2}

	assert.True(t, slices.Equal(expected, result.AllUsers()))
}

func Test_GetNextBirthdayUsers_1(t *testing.T) {
	currentDate := time.Date(time.Now().Year(), 5, 1, 0, 0, 0, 0, time.UTC)

	users := Users{}
	user1 := User{Name: "1"}
	user1.SetBirthday(time.Date(1991, 1, 4, 0, 0, 0, 0, time.UTC), currentDate)
	user2 := User{Name: "2"}
	user2.SetBirthday(time.Date(1991, 2, 13, 0, 0, 0, 0, time.UTC), currentDate)
	user3 := User{Name: "3"}
	user3.SetBirthday(time.Date(1991, 5, 7, 0, 0, 0, 0, time.UTC), currentDate)
	user4 := User{Name: "4"}
	user4.SetBirthday(time.Date(1991, 7, 11, 0, 0, 0, 0, time.UTC), currentDate)
	user5 := User{Name: "5"}
	user5.SetBirthday(time.Date(1991, 12, 2, 0, 0, 0, 0, time.UTC), currentDate)
	users.Add(&user1)
	users.Add(&user2)
	users.Add(&user3)
	users.Add(&user4)
	users.Add(&user5)

	expected := []*User{&user3, &user4, &user5}

	result, _ := users.GetNextBirthdayUsers(3)

	assert.True(t, slices.Equal(expected, result.AllUsers()))
}

func Test_GetNextBirthdayUsers_2(t *testing.T) {
	currentDate := time.Date(time.Now().Year(), 11, 1, 0, 0, 0, 0, time.UTC)

	users := Users{}
	user1 := User{Name: "1"}
	user1.SetBirthday(time.Date(1991, 1, 4, 0, 0, 0, 0, time.UTC), currentDate)
	user2 := User{Name: "2"}
	user2.SetBirthday(time.Date(1991, 2, 13, 0, 0, 0, 0, time.UTC), currentDate)
	user3 := User{Name: "3"}
	user3.SetBirthday(time.Date(1991, 5, 7, 0, 0, 0, 0, time.UTC), currentDate)
	user4 := User{Name: "4"}
	user4.SetBirthday(time.Date(1991, 7, 11, 0, 0, 0, 0, time.UTC), currentDate)
	user5 := User{Name: "5"}
	user5.SetBirthday(time.Date(1991, 12, 2, 0, 0, 0, 0, time.UTC), currentDate)
	users.Add(&user1)
	users.Add(&user2)
	users.Add(&user3)
	users.Add(&user4)
	users.Add(&user5)

	expected := []*User{&user5, &user1, &user2, &user3}

	result, _ := users.GetNextBirthdayUsers(4)

	assert.True(t, slices.Equal(expected, result.AllUsers()))
}

func Test_GetNextBirthdayUsers_3(t *testing.T) {
	currentDate := time.Date(time.Now().Year(), 5, 1, 0, 0, 0, 0, time.UTC)

	users := Users{}
	user1 := User{Name: "1"}
	user1.SetBirthday(time.Date(1991, 1, 4, 0, 0, 0, 0, time.UTC), currentDate)
	user2 := User{Name: "2"}
	user2.SetBirthday(time.Date(1991, 2, 13, 0, 0, 0, 0, time.UTC), currentDate)
	user3 := User{Name: "3"}
	user3.SetBirthday(time.Date(1991, 5, 7, 0, 0, 0, 0, time.UTC), currentDate)
	user4 := User{Name: "4"}
	user4.SetBirthday(time.Date(1991, 7, 11, 0, 0, 0, 0, time.UTC), currentDate)
	user5 := User{Name: "5"}
	user5.SetBirthday(time.Date(1991, 12, 2, 0, 0, 0, 0, time.UTC), currentDate)
	users.Add(&user1)
	users.Add(&user2)
	users.Add(&user3)
	users.Add(&user4)
	users.Add(&user5)

	expected := []*User{&user3, &user4, &user5, &user1, &user2}

	result, _ := users.GetNextBirthdayUsers(999)

	assert.True(t, slices.Equal(expected, result.AllUsers()))
}
