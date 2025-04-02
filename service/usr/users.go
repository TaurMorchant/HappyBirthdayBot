package usr

import (
	"sort"
	"unicode/utf8"
)

type Users struct {
	users []*User
}

//---------------------------------------------------------------------------------------------------------------------

func (u *Users) AllUsers() []*User {
	return u.users
}

func (u *Users) Add(user *User) {
	u.users = append(u.users, user)
	u.sort()
}

func (u *Users) Get(id UserId) (*User, bool) {
	for _, user := range u.users {
		if user.Id == id {
			return user, true
		}
	}
	return nil, false
}

func (u *Users) Delete(id UserId) {
	for i, user := range u.users {
		if user.Id == id {
			u.users = append(u.users[:i], u.users[i+1:]...)
		}
	}
}

func (u *Users) GetMaxNameLength() int {
	var result = 0
	for _, user := range u.users {
		if utf8.RuneCountInString(user.Name) > result {
			result = utf8.RuneCountInString(user.Name)
		}
	}
	return result
}

func (u *Users) GetMaxMonthLength() int {
	var result = 0
	for _, user := range u.users {
		if utf8.RuneCountInString(user.BirthDay().MonthName()) > result {
			result = utf8.RuneCountInString(user.BirthDay().MonthName())
		}
	}
	return result
}

func (u *Users) GetNextBirthdayUsers(n int) (*Users, error) {
	if n > len(u.AllUsers()) {
		n = len(u.AllUsers())
	}

	var users []*User

	sortedUsers := u.sortByDaysBeforeBirthday()

	for i := 0; i < n; i++ {
		user := (sortedUsers.AllUsers())[i]
		users = append(users, user)
	}

	result := &Users{users: users}
	return result, nil
}

//---------------------------------------------------------------------------------------------------------------------

func (u *Users) sort() {
	sort.Slice(u.users, func(i, j int) bool {
		return u.users[i].birthDay.CurrentYearBirthday().Before(u.users[j].birthDay.CurrentYearBirthday())
	})
}

func (u *Users) sortByDaysBeforeBirthday() *Users {
	result := Users{}
	usersClone := make([]*User, len(u.users))
	copy(usersClone, u.users)
	result.users = usersClone

	sort.Slice(result.users, func(i, j int) bool {
		return result.users[i].daysBeforeBirthday < (result.users[j].daysBeforeBirthday)
	})

	return &result
}
