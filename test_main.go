package main

import (
	"happy-birthday-bot/sheets"
	"log"
)

func main() {
	users := sheets.Read()

	log.Println("original usr:", users)

	users.GetAllUsers()[0].Name = users.GetAllUsers()[0].Name + "1"

	log.Println("changed usr:", users)

	sheets.Write(&users)
}
