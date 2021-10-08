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
const agileStartNoteTemplateName = "agile-start-template.md"
const agileStartNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + agileStartNoteTemplateName

func main() {
	optionVal := flag.String("t", "", "テンプレートを指定します。")
	flag.Parse()

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

	data := map[string]interface{}{
		"CreatedAt":            createdAt,
		"Tags":                 tags,
		"Yesterday":            yesterday,
		"Today":                today,
		"Tomorrow":             tomorrow,
		"LastWeeklyFileName":   lastWeeklyFileName,
		"WeeklyFileName":       weeklyFileName,
		"NextWeeklyFileName":   nextWeeklyFileName,
		"DailyNote1":           dailyNote1,
		"DailyNote2":           dailyNote2,
		"DailyNote3":           dailyNote3,
		"DailyNote4":           dailyNote4,
		"DailyNote5":           dailyNote5,
		"DailyNote6":           dailyNote6,
		"DailyNote7":           dailyNote7,
		"LastWeekGoals":        extractLastWeekGoals(getLastWeeklyReview()),
		"RemainingTasks":       ExtractYesterdayTasks(dailyNotePath, yesterday),
		"WeeklyFDL":            ExtractWeeklyFDL(dailyNotePath, []string{dailyNote1, dailyNote2, dailyNote3, dailyNote4, dailyNote5, dailyNote6}),
		"zettelkastenFileName": zettelkastenFileName,
	}

	var filePath string
	switch *optionVal {
	case "daily":
		createDailyNote(data, dailyNotePath)
		filePath = dailyNotePath + data["Today"].(string) + ".md"
	case "weekly":
		createWeeklyNote(data, weeklyNotePath, t.Weekday().String())
		filePath = weeklyNotePath + data["WeeklyFileName"].(string) + ".md"
	case "zettelkasten":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"].(string), zettelkastenNoteTemplatePath, zettelkastenNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"].(string) + ".md"
	case "research":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"].(string), researchLogNoteTemplatePath, researchLogNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"].(string) + ".md"
	case "agile-start":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"].(string), agileStartNoteTemplatePath, agileStartNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"].(string) + ".md"
	default:
		fmt.Printf("テンプレートを正しく指定してください %s\n", *optionVal)
	}
	openCreatedFile(filePath)
}

func createDailyNote(data map[string]interface{}, dailyNotePath string) {
	// 本日日付のファイルがなければ作成
	if _, err := os.Stat(dailyNotePath + data["Today"].(string) + ".md"); os.IsNotExist(err) {
		createNoteFromTemplate(data, dailyNotePath, data["Today"].(string), dailyNoteTemplatePath, dailyNoteTemplateName)
	} else {
		log.Println("Daily notesはすでに存在します。")
	}
}

func createWeeklyNote(data map[string]interface{}, weeklyNotePath string, weekDay string) {
	// 土曜日の場合かつ該当するファイルがなければweekly-reviewファイルを作成
	if _, err := os.Stat(weeklyNotePath + data["WeeklyFileName"].(string) + ".md"); os.IsNotExist(err) && weekDay == "Friday" {
		createNoteFromTemplate(data, weeklyNotePath, data["WeeklyFileName"].(string), weeklyNoteTemplatePath, weeklyNoteTemplateName)
	} else {
		log.Println("Weekly reviewsはすでに存在します。")
	}
}

func createNoteFromTemplate(data map[string]interface{}, notePath string, noteFileName string, templatePath string, templateName string) {
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

func ExtractYesterdayTasks(dailyNotePath string, dailyFileName string) string {
	text, err := ioutil.ReadFile(dailyNotePath + dailyFileName + ".md")
	if err != nil {
		log.Fatal(err)
		return ""
	}
	r := regexp.MustCompile(`(?m)^## TODO.*[\s\S\n]*?\n$`)
	// remove "## TODO"
	result := strings.Replace(r.FindString(string(text)), "## TODO\n", "", 1)

	return strings.TrimRight(result, "\n")
}

func ExtractWeeklyFDL(dailyNotePath string, dailyFileNames []string) map[string]string {
	result := map[string]string{}
	for _, v := range dailyFileNames {
		text, err := ioutil.ReadFile(dailyNotePath + v + ".md")
		if err != nil {
			log.Fatal(err)
			return result
		}
		// extract Fun
		r1 := regexp.MustCompile(`(?m)^- Fun.*[\s\S\n]*^- Done`)
		rep1 := regexp.MustCompile(`(?m)- Fun|- Done`)
		fun := rep1.ReplaceAllString(r1.FindString(string(text)), "")
		result["Fun"] += strings.Trim(fun, "\n")

		// extract Done
		r2 := regexp.MustCompile(`(?m)^- Done.*[\s\S\n]*^- Learn`)
		rep2 := regexp.MustCompile(`(?m)- Done|- Learn`)
		done := rep2.ReplaceAllString(r2.FindString(string(text)), "")
		result["Done"] += strings.Trim(done, "\n")

		// extract Learn
		r3 := regexp.MustCompile(`(?m)^- Learn.*[\s\S\n]*?\n$`)
		rep3 := regexp.MustCompile(`(?m)- Learn`)
		learn := rep3.ReplaceAllString(r3.FindString(string(text)), "")
		result["Learn"] += strings.Trim(learn, "\n")
	}
	return result
}

func extractLastWeekGoals(lastWeeklyFileName string) string {
	text, err := ioutil.ReadFile(weeklyNotePath + lastWeeklyFileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	r := regexp.MustCompile(`(?m)^- Try.*[\s\S\n]*?\n$`)
	// remove "- Try"
	result := strings.Replace(r.FindString(string(text)), "- Try\n", "", 1)

	// trim head spaces
	r2 := regexp.MustCompile(`(?m)^\s{2}`)
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
	_, _ = cmd.Output()
}
