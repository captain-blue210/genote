package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const dailyNotePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/look-back/daily-notes/"
const dailyNoteTemplateName = "daily-note-template-v2.md"
const dailyNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + dailyNoteTemplateName

const weeklyNotePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/look-back/weekly-reviews/"
const weeklyNoteTemplateName = "weekly-review-template-v2.md"
const weeklyNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + weeklyNoteTemplateName

const notePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/notes/"
const zettelkastenNoteTemplateName = "zettelkasten-template.md"
const zettelkastenNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + zettelkastenNoteTemplateName
const researchLogNoteTemplateName = "research-log-template.md"
const researchLogNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + researchLogNoteTemplateName

func main() {
	optionVal := flag.String("t", "", "テンプレートを指定します。")
	flag.Parse()

	t := time.Now()
	data := map[string]string{
		"CreatedAt":            t.Format("2006-01-02 15:04:05"),
		"Tags":                 strconv.Itoa(t.Year()) + "/" + t.Format("01"),
		"Yesterday":            t.AddDate(0, 0, -1).Format("2006-01-02"),
		"Today":                t.Format("2006-01-02"),
		"Tomorrow":             t.AddDate(0, 0, 1).Format("2006-01-02"),
		"LastWeeklyFileName":   t.AddDate(0, 0, -13).Format("20060102") + "-" + t.AddDate(0, 0, -7).Format("20060102"),
		"WeeklyFileName":       t.AddDate(0, 0, -6).Format("20060102") + "-" + t.Format("20060102"),
		"NextWeeklyFileName":   t.AddDate(0, 0, 1).Format("20060102") + "-" + t.AddDate(0, 0, 7).Format("20060102"),
		"DailyNote1":           t.AddDate(0, 0, -6).Format("2006-01-02"),
		"DailyNote2":           t.AddDate(0, 0, -5).Format("2006-01-02"),
		"DailyNote3":           t.AddDate(0, 0, -4).Format("2006-01-02"),
		"DailyNote4":           t.AddDate(0, 0, -3).Format("2006-01-02"),
		"DailyNote5":           t.AddDate(0, 0, -2).Format("2006-01-02"),
		"DailyNote6":           t.AddDate(0, 0, -1).Format("2006-01-02"),
		"DailyNote7":           t.Format("2006-01-02"),
		"LastWeekGoals":        extractLastWeekGoals(getLastWeeklyReview()),
		"zettelkastenFileName": t.Format("2006-01-02-15-04-05"),
	}

	var filePath string
	switch *optionVal {
	case "daily":
		createDailyNote(data, dailyNotePath)
		filePath = dailyNotePath + data["Today"] + ".md"
	case "weekly":
		createWeeklyNote(data, weeklyNotePath, t.Weekday().String())
		filePath = weeklyNotePath + data["WeeklyFileName"] + ".md"
	case "zettelkasten":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"], zettelkastenNoteTemplatePath, zettelkastenNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"] + ".md"
	case "research":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"], researchLogNoteTemplatePath, researchLogNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"] + ".md"
	default:
		fmt.Printf("テンプレートを正しく指定してください %s\n", *optionVal)
	}

	openCreatedFile(filePath)
}

func createDailyNote(data map[string]string, dailyNotePath string) {
	// 本日日付のファイルがなければ作成
	if _, err := os.Stat(dailyNotePath + data["Today"] + ".md"); os.IsNotExist(err) {
		createNoteFromTemplate(data, dailyNotePath, data["Today"], dailyNoteTemplatePath, dailyNoteTemplateName)
	} else {
		log.Println("Daily notesはすでに存在します。")
	}
}

func createWeeklyNote(data map[string]string, weeklyNotePath string, weekDay string) {
	// 土曜日の場合かつ該当するファイルがなければweekly-reviewファイルを作成
	if _, err := os.Stat(weeklyNotePath + data["WeeklyFileName"] + ".md"); os.IsNotExist(err) && weekDay == "Saturday" {
		createNoteFromTemplate(data, weeklyNotePath, data["WeeklyFileName"], weeklyNoteTemplatePath, weeklyNoteTemplateName)
	} else {
		log.Println("Weekly reviewsはすでに存在します。")
	}
}

func createNoteFromTemplate(data map[string]string, notePath string, noteFileName string, templatePath string, templateName string) {
	f, err := os.Create(notePath + noteFileName + ".md")
	if err != nil {
		log.Fatal(err)
	}
	te, err := template.New(templateName).ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
	}
	if err = te.Execute(f, data); err != nil {
		log.Fatal(err)
	}
}

func extractLastWeekGoals(lastWeeklyFileName string) string {
	text, err := ioutil.ReadFile(weeklyNotePath + lastWeeklyFileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	r := regexp.MustCompile(`(?m)^- Try.*[\s\S\n]*?\n$`)
	// remove "- Try"
	result := strings.Replace(r.FindString(string(text)), "- Try", "", 1)

	// trim head spaces
	r2 := regexp.MustCompile(`(?m)^\s`)
	result = r2.ReplaceAllString(result, "")
	return strings.TrimRight(result, "\n")
}

func getLastWeeklyReview() string {
	files, err := ioutil.ReadDir(weeklyNotePath)
	if err != nil {
		log.Fatal(err)
	}

	var fileNames = []time.Time{}
	for _, v := range files {
		t, _ := time.Parse("20060102", strings.Replace(strings.Split(v.Name(), "-")[1], ".md", "", 1))
		fileNames = append(fileNames, t)
	}

	sort.Slice(fileNames, func(i, j int) bool {
		return fileNames[i].After(fileNames[j])
	})

	var result string
	for _, v := range files {
		if fileNames[0].Format("20060102") == strings.Replace(strings.Split(v.Name(), "-")[1], ".md", "", 1) {
			result = v.Name()
		}
	}
	return result
}

func openCreatedFile(path string) {
	cmd := exec.Command("code", "-r", path)
	result, _ := cmd.Output()
	log.Println(string(result))
}
