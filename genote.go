package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const DailyNotePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/look-back/daily-notes/"
const dailyNoteTemplateName = "daily-note-template-v2.md"
const dailyNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + dailyNoteTemplateName

const weeklyNotePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/look-back/weekly-reviews/"
const weeklyNoteTemplateName = "weekly-review-template-v2.md"
const weeklyNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + weeklyNoteTemplateName

const MonthlyNotePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/look-back/monthly-reviews/"
const monthlyNoteTemplateName = "monthly-review-template.md"
const monthlyNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + monthlyNoteTemplateName

const YearlyNotePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/look-back/yearly-reviews/"

const notePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/notes/"
const zettelkastenNoteTemplateName = "zettelkasten-template.md"
const zettelkastenNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + zettelkastenNoteTemplateName
const researchLogNoteTemplateName = "research-log-template.md"
const researchLogNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + researchLogNoteTemplateName
const agileStartNoteTemplateName = "agile-start-template.md"
const agileStartNoteTemplatePath = "/Users/captain-blue/Library/Mobile Documents/iCloud~md~obsidian/Documents/second-brain/templates/" + agileStartNoteTemplateName

func main() {
	optionVal := flag.String("t", "", "テンプレートを指定します。")
	datetimeVal := flag.String("d", "", "日付を指定します。")
	backDaysVal := flag.Int("bd", 6, "メモ出力でさかのぼる日数を指定します")
	flag.Parse()

	t := time.Now()

	if *datetimeVal != "" {
		t, _ = time.Parse("2006-01-02", *datetimeVal)
	}

	data := CreateBasicData(datetimeVal, t)

	var filePath string
	switch *optionVal {
	case "daily":
		createDailyNote(data, DailyNotePath)
		filePath = DailyNotePath + data["Today"].(string) + ".md"
	case "weekly":
		createWeeklyNote(data, weeklyNotePath, t.Weekday().String())
		filePath = weeklyNotePath + data["WeeklyFileName"].(string) + ".md"
	case "monthly":
		CreateMonthlyNote(data, YearlyNotePath, MonthlyNotePath, weeklyNotePath, t)
		filePath = MonthlyNotePath + data["MonthlyFileName"].(string) + ".md"
	case "zettelkasten":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"].(string), zettelkastenNoteTemplatePath, zettelkastenNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"].(string) + ".md"
	case "research":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"].(string), researchLogNoteTemplatePath, researchLogNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"].(string) + ".md"
	case "agile-start":
		createNoteFromTemplate(data, notePath, data["zettelkastenFileName"].(string), agileStartNoteTemplatePath, agileStartNoteTemplateName)
		filePath = notePath + data["zettelkastenFileName"].(string) + ".md"
	case "memo":
		fmt.Println(ExtractMemo(DailyNotePath, *backDaysVal, t))
		os.Exit(0)
	default:
		fmt.Printf("テンプレートを正しく指定してください %s\n", *optionVal)
		os.Exit(1)
	}
	openCreatedFile(filePath)
}

func CreateBasicData(datetimeVal *string, t time.Time) map[string]interface{} {
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
	lastMonthlyFileName := t.AddDate(0, -1, 0).Format("2006-01")
	monthlyFileName := t.Format("2006-01")
	nextMonthlyFileName := t.AddDate(0, 1, 0).Format("2006-01")

	return map[string]interface{}{
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
		"RemainingTasks":       ExtractYesterdayTasks(DailyNotePath, yesterday),
		"WeeklyFDL":            ExtractWeeklyFDL(DailyNotePath, []string{dailyNote1, dailyNote2, dailyNote3, dailyNote4, dailyNote5, dailyNote6}),
		"zettelkastenFileName": zettelkastenFileName,
		"LastMonthlyFileName":  lastMonthlyFileName,
		"MonthlyFileName":      monthlyFileName,
		"NextMonthlyFileName":  nextMonthlyFileName,
	}
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
	if _, err := os.Stat(weeklyNotePath + data["WeeklyFileName"].(string) + ".md"); os.IsNotExist(err) {
		createNoteFromTemplate(data, weeklyNotePath, data["WeeklyFileName"].(string), weeklyNoteTemplatePath, weeklyNoteTemplateName)
	} else {
		log.Println("Weekly reviewsはすでに存在します。")
	}
}

func CreateMonthlyNote(data map[string]interface{}, yearlyNotePath string, monthlyNotePath string, weeklyNotePath string, currentTime time.Time) {
	data["Period"] = GetCurrentQuarter(currentTime.Month())
	data["OKR"] = GetCurrentOKR(yearlyNotePath, strconv.Itoa(currentTime.Year()), GetCurrentQuarter(currentTime.Month()))
	data["MonthlyKPT"] = ExtractMonthlyKPT(weeklyNotePath, currentTime)
	data["MonthlyTechArticles"] = ExtractMonthlyArticles(currentTime)

	if _, err := os.Stat(monthlyNotePath + data["MonthlyFileName"].(string) + ".md"); os.IsNotExist(err) {
		createNoteFromTemplate(data, monthlyNotePath, data["MonthlyFileName"].(string), monthlyNoteTemplatePath, monthlyNoteTemplateName)
	} else {
		log.Println("Monthly reviewはすでに存在します。")
	}
}

func ExtractMonthlyKPT(weeklyNotePath string, currentTime time.Time) map[string]string {
	result := map[string]string{}
	/* weekly-reviewsからx月のものを抜き出す
	週次レビューが月をまたいでいる場合も対象とする
	*/
	targetStr := currentTime.Format("200601")
	files, _ := ioutil.ReadDir(weeklyNotePath)
	for i, f := range files {
		if !strings.Contains(f.Name(), targetStr) {
			continue
		}

		text, err := ioutil.ReadFile(weeklyNotePath + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		// 上記で取得したファイルからKeep, Problem, Tryを抜き出しMapに保存する
		// extract Keep
		r1 := regexp.MustCompile(`(?m)^- Keep.*[\s\S\n]*^- Problem`)
		rep1 := regexp.MustCompile(`(?m)- Keep|- Problem`)
		keep := rep1.ReplaceAllString(r1.FindString(string(text)), "")
		result["Keep"] += trim(i, keep)

		// extract Probrem
		r2 := regexp.MustCompile(`(?m)^- Problem.*[\s\S\n]*^- Try`)
		rep2 := regexp.MustCompile(`(?m)- Problem|- Try`)
		probrem := rep2.ReplaceAllString(r2.FindString(string(text)), "")
		result["Probrem"] += trim(i, probrem)

		// extract Try
		r3 := regexp.MustCompile(`(?m)^- Try.*[\s\S\n]*?\n$`)
		rep3 := regexp.MustCompile(`(?m)- Try`)
		try := rep3.ReplaceAllString(r3.FindString(string(text)), "")
		result["Try"] += trim(i, try)
	}
	return result
}

func GetCurrentQuarter(currentMonth time.Month) string {
	var result string
	switch currentMonth {
	case time.January, time.February, time.March:
		result = "1~3月"
	case time.April, time.May, time.June:
		result = "4~6月"
	case time.July, time.August, time.September:
		result = "7~9月"
	case time.October, time.November, time.December:
		result = "10~12月"
	}
	return result
}

func GetCurrentOKR(yearlyNotePath string, currentYear string, currentQuarter string) map[string]string {
	result := map[string]string{}
	// currentQuarter ~ 「### ふりかえり」の直前を取得する
	text, err := ioutil.ReadFile(yearlyNotePath + currentYear + ".md")
	if err != nil {
		log.Fatal(err)
	}

	extractRegex := regexp.MustCompile(`(?m)^###\s` + currentQuarter + `.*[\s\S\n]*## ふりかえり`)
	removeRegex := regexp.MustCompile(`## ふりかえり|### ` + currentQuarter)

	// GoalsAndResultsをキーにして取得した文字列をバリューに入れる
	result["GoalsAndResults"] = strings.Trim(removeRegex.ReplaceAllString(extractRegex.FindString(string(text)), ""), "\n")

	return result
}

type AutoGenerated struct {
	Articles []struct {
		Title       string    `json:"title"`
		Slug        string    `json:"slug"`
		PublishedAt time.Time `json:"published_at"`
	} `json:"articles"`
}

func ExtractMonthlyArticles(target time.Time) string {
	const BASE_URL = "https://zenn.dev/captain_blue/articles/"
	// ZennのAPIを叩いてJSONを取得する
	raw := execAPI("https://zenn.dev/api/articles?username=captain_blue&count=500&order=latest")
	parsed := AutoGenerated{}
	err := json.Unmarshal([]byte(raw), &parsed)
	if err != nil {
		log.Println(err)
	}
	// 取得したJSONのpublished_atとtargetを比較し、同一の月のみに絞り込む
	var result string
	for _, v := range parsed.Articles {
		if target.Year() == v.PublishedAt.Year() && target.Month() == v.PublishedAt.Month() {
			// 該当するデータのslugを取得し、BASE_URLと結合する
			result += v.Title + "\n" + BASE_URL + v.Slug + "\n"
		}
	}

	return result
}

func execAPI(url string) string {
	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(body)
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
	r := regexp.MustCompile(`(?m)^## TODO.*[\s\S\n]*## ふりかえり`)
	rep := regexp.MustCompile(`(?m)## TODO\n\n|## ふりかえり`)
	// remove "## TODO"
	result := rep.ReplaceAllString(r.FindString(string(text)), "")

	return strings.Trim(result, "\n")
}

func ExtractWeeklyFDL(dailyNotePath string, dailyFileNames []string) map[string]string {
	result := map[string]string{}
	for i, v := range dailyFileNames {
		text, err := ioutil.ReadFile(dailyNotePath + v + ".md")
		if err != nil {
			log.Fatal(err)
			return result
		}
		// extract Fun
		r1 := regexp.MustCompile(`(?m)^- Fun.*[\s\S\n]*^- Done`)
		rep1 := regexp.MustCompile(`(?m)- Fun|- Done`)
		fun := rep1.ReplaceAllString(r1.FindString(string(text)), "")
		result["Fun"] += trim(i, fun)

		// extract Done
		r2 := regexp.MustCompile(`(?m)^- Done.*[\s\S\n]*^- Learn`)
		rep2 := regexp.MustCompile(`(?m)- Done|- Learn`)
		done := rep2.ReplaceAllString(r2.FindString(string(text)), "")
		result["Done"] += trim(i, done)
		// extract Learn
		r3 := regexp.MustCompile(`(?m)^- Learn.*[\s\S\n]*?\n$`)
		rep3 := regexp.MustCompile(`(?m)- Learn`)
		learn := rep3.ReplaceAllString(r3.FindString(string(text)), "")
		result["Learn"] += trim(i, learn)
	}
	return result
}

func trim(count int, text string) string {
	if count == 0 {
		// 1回目は文頭、文末の改行を削除
		return strings.Trim(text, "\n")
	} else {
		// 2回目以降は文末の改行のみ削除
		return strings.TrimRight(text, "\n")
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

func createDailyNotePath(days int, dailyNotePath string, t time.Time) []string {
	var result []string
	for i := -days; i < 0; i++ {
		dailyNote := t.AddDate(0, 0, i).Format("2006-01-02")
		path := dailyNotePath + dailyNote + ".md"
		result = append(result, path)
	}
	return result
}

func ExtractMemo(dailyNotePath string, days int, t time.Time) string {
	result := "\n"
	paths := createDailyNotePath(days, dailyNotePath, t)
	for i, f := range paths {
		text, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
			return result
		}

		// extract メモ
		r1 := regexp.MustCompile(`(?m)^## メモ.*[\s\S\n]*`)
		rep1 := regexp.MustCompile(`(?m)## メモ`)
		memo := rep1.ReplaceAllString(r1.FindString(string(text)), "")
		result += trim(i, memo)
	}
	return result
}
