package main

import (
	"fmt"
	"time"
)

const (
	studyType     = "Бакалавриат"
	timePeriods   = 7
	botKey        = "6681335182:AAGCpubq0i_jizuUdlqQbZY5OFTbdtWn5QA"
	channelChatId = -1002162184959
)

func main() {
	groupNumber := "11-402"
	filename := "Расписание консультаций ИТИС КФУ с 01.09.2024 по 17.10.2024.xlsx"
	groups := ParseFile(filename)

	// setup bot
	bot := CreateBot(botKey)

	// setting up reminders for specific group
	notifier := CreateNotifier(withTGBot(bot))
	notifier.SetupLessonsReminders(groups[groupNumber])

	// chillin
	fmt.Println("Waiting...")
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
