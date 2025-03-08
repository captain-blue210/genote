package main

import (
	"strconv"
	"time"
)

type NoteData struct {
	CreatedAt            string
	Tags                 string
	Yesterday            string
	Today                string
	Tomorrow             string
	LastWeeklyFileName   string
	WeeklyFileName       string
	NextWeeklyFileName   string
	DailyNote1           string
	DailyNote2           string
	DailyNote3           string
	DailyNote4           string
	DailyNote5           string
	DailyNote6           string
	DailyNote7           string
	LastWeekGoals        string
	RemainingTasks       string
	WeeklyFDL            map[string]string
	ZettelkastenFileName string
}

func createNoteData() NoteData {

	t := time.Now()
	createdAt := t.Format("2006-01-02 15:04:05")
	tags := strconv.Itoa(t.Year()) + "/" + t.Format("01")
	yesterday := t.AddDate(0, 0, -1).Format("2006-01-02")
	today := t.Format("2006-01-02")
	tomorrow := t.AddDate(0, 0, 1).Format("2006-01-02")
	lastWeeklyFileName := t.AddDate(0, 0, -13).Format("20060102") + "-" + t.AddDate(0, 0, -7).Format("20060102")
	weeklyFileName := t.AddDate(0, 0, -6).Format("20060102") + "-" + t.Format("20060102")
	nextWeeklyFileName := t.AddDate(0, 0, 1).Format("20060102") + "-" + t.AddDate(0, 0, 7).Format("20060102")
	dailyNote1 := t.AddDate(0, 0, -6).Format("2006-01-02")
	dailyNote2 := t.AddDate(0, 0, -5).Format("2006-01-02")
	dailyNote3 := t.AddDate(0, 0, -4).Format("2006-01-02")
	dailyNote4 := t.AddDate(0, 0, -3).Format("2006-01-02")
	dailyNote5 := t.AddDate(0, 0, -2).Format("2006-01-02")
	dailyNote6 := t.AddDate(0, 0, -1).Format("2006-01-02")
	dailyNote7 := t.Format("2006-01-02")
	zettelkastenFileName := t.Format("2006-01-02-15-04-05")

	return NoteData{
		createdAt,
		tags,
		yesterday,
		today,
		tomorrow,
		lastWeeklyFileName,
		weeklyFileName,
		nextWeeklyFileName,
		dailyNote1,
		dailyNote2,
		dailyNote3,
		dailyNote4,
		dailyNote5,
		dailyNote6,
		dailyNote7,
		extractLastWeekGoals(getLastWeeklyReview()),
		ExtractYesterdayTasks(DailyNotePath, yesterday),
		ExtractWeeklyFDL(DailyNotePath, []string{dailyNote1, dailyNote2, dailyNote3, dailyNote4, dailyNote5, dailyNote6}),
		zettelkastenFileName,
	}
}
