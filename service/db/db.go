package db

import (
	"database/sql"
	"happy-birthday-bot/date"
	"happy-birthday-bot/usr"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var instance *sql.DB

func Init(path string) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Panic("Failed to create db directory:", err)
	}

	var err error
	instance, err = sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Panic("Failed to open SQLite database:", err)
	}
	instance.SetMaxOpenConns(1)
	if err = instance.Ping(); err != nil {
		log.Panic("Failed to connect to SQLite database:", err)
	}

	_, err = instance.Exec(`CREATE TABLE IF NOT EXISTS users (
		id                 INTEGER PRIMARY KEY,
		name               TEXT    NOT NULL,
		birthday           TEXT    NOT NULL,
		wishlist           TEXT    NOT NULL DEFAULT '',
		reminder30days     INTEGER NOT NULL DEFAULT 0,
		reminder15days     INTEGER NOT NULL DEFAULT 0,
		birthday_greetings INTEGER NOT NULL DEFAULT 0,
		disabled           INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		log.Panic("Failed to init users table:", err)
	}

	_, err = instance.Exec(`ALTER TABLE users ADD COLUMN disabled INTEGER NOT NULL DEFAULT 0`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		log.Panic("Failed to migrate users table:", err)
	}

	log.Println("SQLite initialized:", path)
}

func ReadUsers() usr.Users {
	rows, err := instance.Query(
		`SELECT id, name, birthday, wishlist, reminder30days, reminder15days, birthday_greetings FROM users WHERE disabled = 0`,
	)
	if err != nil {
		log.Panic("Failed to read users:", err)
	}
	defer rows.Close()

	var users usr.Users
	timeNow := time.Now()
	for rows.Next() {
		var id int64
		var name, birthdayStr, wishlist string
		var r30, r15, greetings int
		if err := rows.Scan(&id, &name, &birthdayStr, &wishlist, &r30, &r15, &greetings); err != nil {
			log.Panic("Failed to scan user row:", err)
		}
		birthday, err := date.ParseBirthday(birthdayStr)
		if err != nil {
			log.Panicf("Invalid birthday %q for user %d in database: %v", birthdayStr, id, err)
		}
		user := &usr.User{
			Id:                usr.UserId(id),
			Name:              name,
			Wishlist:          wishlist,
			Reminder30days:    r30 != 0,
			Reminder15days:    r15 != 0,
			BirthdayGreetings: greetings != 0,
		}
		user.SetBirthday2(birthday, timeNow)
		users.Add(user)
	}
	if err := rows.Err(); err != nil {
		log.Panic("Failed to iterate users:", err)
	}
	return users
}

func InsertUser(user *usr.User) error {
	_, err := instance.Exec(
		`INSERT INTO users (id, name, birthday, wishlist) VALUES (?, ?, ?, ?)`,
		int64(user.Id), user.Name, user.BirthDay().ToString(), user.Wishlist,
	)
	return err
}

func DisableUser(userId usr.UserId) error {
	_, err := instance.Exec(`UPDATE users SET disabled = 1 WHERE id = ?`, int64(userId))
	return err
}

func ReactivateUser(userId usr.UserId, name string, birthday date.Birthday) (bool, error) {
	result, err := instance.Exec(
		`UPDATE users SET disabled = 0, name = ?, birthday = ? WHERE id = ? AND disabled = 1`,
		name, birthday.ToString(), int64(userId),
	)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	return rows > 0, err
}

func UpdateWishlist(userId usr.UserId, wishlist string) error {
	_, err := instance.Exec(`UPDATE users SET wishlist = ? WHERE id = ?`, wishlist, int64(userId))
	return err
}

func UpdateFlags(user *usr.User) error {
	_, err := instance.Exec(
		`UPDATE users SET reminder30days = ?, reminder15days = ?, birthday_greetings = ? WHERE id = ?`,
		boolToInt(user.Reminder30days), boolToInt(user.Reminder15days), boolToInt(user.BirthdayGreetings),
		int64(user.Id),
	)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
