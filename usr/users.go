package usr

import (
	"fmt"
	"happy-birthday-bot/date"
	"sort"
	"unicode/utf8"
)

type UserId int64

type User struct {
	Id       UserId
	Name     string
	Birthday date.Birthday
}

func (u User) FormattedString(maxNameLength int) string {
	return fmt.Sprintf("🎂 %*s — %-10s", maxNameLength, u.Name, u.Birthday.ToString())
}

//------------------

type Users struct {
	users []User
}

func (u *Users) GetAllUsers() []User {
	return u.users
}

func (u *Users) Add(user User) {
	u.users = append(u.users, user)
	u.sort()
}

func (u *Users) Get(id UserId) (*User, bool) {
	for _, user := range u.users {
		if user.Id == id {
			return &user, true
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

//-----------------------------------------

func (u *Users) sort() {
	sort.Slice(u.users, func(i, j int) bool {
		return u.users[i].Birthday.Time.Before(u.users[j].Birthday.Time)
	})
}
