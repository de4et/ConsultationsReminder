package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/uniplaces/carbon"
)

type Notifier struct {
	bot *TGBot
}

var levelConverter = map[int]int{
	1: 30,
	2: 15,
	3: 5,
}

const template = `
%s

Через %d минут начнётся:
%s
`

func CreateNotifier(opts ...func(*Notifier)) *Notifier {
	ntfr := &Notifier{}
	for _, o := range opts {
		o(ntfr)
	}
	return ntfr
}

func (n *Notifier) SetupLessonsReminders(group *Group) {
	for _, v := range group.Lessons {
		setupLessonReminder(n.bot, v)
	}
}

func getMinutesLeftByLevel(level int) int {
	return levelConverter[level]
}

func formatLessonTemplateMessage(l Lesson, level int) string {
	return fmt.Sprintf(template, strings.Repeat("❗", level), getMinutesLeftByLevel(level), l.Name)
}

func getDiffToNextFreeLesson(l Lesson) time.Duration {
	now := carbon.Now()
	nextDate, _ := carbon.Create(now.Year(), now.Month(), now.Day(), l.Start.Hour, l.Start.Minute, 0, 0, now.TimeZone())
	if now.Weekday() != time.Weekday(l.Day.num) {
		nextDate = now.Next(time.Weekday(l.Day.num))
		nextDate.SetHour(l.Start.Hour)
		nextDate.SetMinute(l.Start.Minute)
	}

	diff := now.DiffInSeconds(nextDate, false)
	if diff < 0 {
		nextDate = now.Next(time.Weekday(l.Day.num))
		nextDate.SetHour(l.Start.Hour)
		nextDate.SetMinute(l.Start.Minute)
		diff = now.DiffInSeconds(nextDate, false)
	}
	return time.Duration(diff) * time.Second
}

func setTimerByLevel(bot *TGBot, l Lesson, diff time.Duration, level int) {
	if diff < 0 {
		return
	}

	time.AfterFunc(diff, func() {
		slog.Info("sending mess", slog.String("level", fmt.Sprintf("%v", level)))
		bot.SendNotification(formatLessonTemplateMessage(l, level))
	})
}

func setupLessonReminder(bot *TGBot, l Lesson) {
	diff := getDiffToNextFreeLesson(l)
	setTimerByLevel(bot, l, diff-30*time.Minute, 1)
	setTimerByLevel(bot, l, diff-15*time.Minute, 2)
	setTimerByLevel(bot, l, diff-5*time.Minute, 3)
}

func withTGBot(bot *TGBot) func(*Notifier) {
	return func(n *Notifier) {
		n.bot = bot
	}
}
