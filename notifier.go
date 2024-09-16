package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/uniplaces/carbon"
)

type NodificationSender interface {
	SendNotification(text string) error
}

type Notifier struct {
	ns NodificationSender
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

func CreateNotifier(ns NodificationSender) *Notifier {
	ntfr := &Notifier{
		ns: ns,
	}
	return ntfr
}

func (n *Notifier) SetupLessonsReminders(group *Group, weeks int) {
	for _, v := range group.Lessons {
		n.setupLessonReminders(v, weeks)
	}
}

func getMinutesLeftByLevel(level int) int {
	return levelConverter[level]
}

func formatLessonTemplateMessage(l Lesson, level int) string {
	return fmt.Sprintf(template, strings.Repeat("❗", level), getMinutesLeftByLevel(level), l.Name)
}

func getDiffToNextFreeLesson(l Lesson, now *carbon.Carbon) time.Duration {
	nextDate, _ := carbon.Create(now.Year(), now.Month(), now.Day(), l.Start.Hour, l.Start.Minute, 0, 0, now.TimeZone())
	if now.Weekday() != time.Weekday(l.Day.num) {
		nextDate = now.Next(time.Weekday(l.Day.num))
		nextDate.SetHour(l.Start.Hour)
		nextDate.SetMinute(l.Start.Minute)
	}

	nowNow := carbon.Now()
	diff := nowNow.DiffInSeconds(nextDate, false)
	if diff < 0 {
		nextDate = now.Next(time.Weekday(l.Day.num))
		nextDate.SetHour(l.Start.Hour)
		nextDate.SetMinute(l.Start.Minute)
		diff = nowNow.DiffInSeconds(nextDate, false)
	}
	return time.Duration(diff) * time.Second
}

func (n *Notifier) setTimerByLevel(l Lesson, diff time.Duration, level int) {
	if diff < 0 {
		return
	}

	time.AfterFunc(diff, func() {
		slog.Info("sending mess", slog.String("level", fmt.Sprintf("%v", level)))
		n.ns.SendNotification(formatLessonTemplateMessage(l, level))
	})
}

func (n *Notifier) setupLessonReminders(l Lesson, weeks int) {
	now := carbon.Now()
	for i := 0; i < weeks; i++ {
		diff := getDiffToNextFreeLesson(l, now.AddWeeks(i))
		n.setTimerByLevel(l, diff-30*time.Minute, 1)
		n.setTimerByLevel(l, diff-15*time.Minute, 2)
		n.setTimerByLevel(l, diff-5*time.Minute, 3)
	}
}
