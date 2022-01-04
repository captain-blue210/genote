package main

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func TestExtractYesterdayTasks(t *testing.T) {
	expected := `- 仕事
   - [x] AUDIT質問回答
   - [x] チケット対応
 - [ ] 健康診断日決める
 - プライベート
   - [ ] 英文解釈教室`

	result := ExtractYesterdayTasks("./test-file/", "test1")
	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:\n" + result)
	} else {
		t.Log("TestExtractYesterdayTasks passed : \n", result)
	}
}

func TestExtractWeeklyFDL(t *testing.T) {
	expected := map[string]string{
		"Fun":   "  - 商品詳細が表示できない問題を解決できた\n  - オペレーション改善ツールについて話ができた\n  - 明日は代休をとれた\n  - 英文解釈教室、伝わる英語表現法、シェルワンライナーを進められた\n  - 家事などをしっかりこなせて、ベランダ掃除もできた",
		"Done":  "  - AUDIT回答",
		"Learn": "  - 最近4~5時間しか眠れていないことがfitbitの記録でわかった\n  - 同格のthat節、接続詞としてのevery time等\n  - ファイルの一括変換、リネーム、xargsの並列実行法、`time`の使い方\n  - 日本語を情報単位で分割して訳す方法",
	}

	fileNames := []string{"test1", "test2"}
	result := ExtractWeeklyFDL("test-file/", fileNames)
	for key := range result {
		if result[key] != expected[key] {
			log.Println("result : " + result[key])
			log.Println("expected : " + expected[key])
			t.Error("抽出した文字列が想定と異なります")
		} else {
			t.Log("TestExtractWeeklyFDL passed : \n", result)
		}
	}
}

func TestExtractLastWeekGoals(t *testing.T) {
	expected := "- 週一発信を再開する\n- Flutterで生活費精算アプリ作成\n- 英文解釈教室を毎日１例題\n- 業務で楽しかったこと、楽しくなかったことを意識的に記録する"
	result := extractLastWeekGoals("20210926-20211002.md")
	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:\n" + result)
	} else {
		t.Log("TestExtractYesterdayTasks passed : \n", result)
	}
}

func TestExtractMemo(t *testing.T) {
	expected := `
- インセプションデッキ
  - JuJu
  - 美味しく焼き肉を食べたい

- 今日のins-pre環境整備で、どうしてもVue側でLOCALのenvファイル内容が使われてしまう問題にハマっていた
  - ローカルのhostsファイルにins-preのドメインを登録していたのを忘れていた
  - どうしたらもっと早く気付けたか？
    - ins-preのitem-webにはちゃんと.env_ins-preの内容が使われていた
    - ここで、hostsファイルを思いつければよかった
    - k8sの設定ミスなどを調べようとしていたのがよくなかった
    - 「.envには正しい内容が登録されているから、他にenvがLOCALになる原因はなにか？」を考えられればよかった
    - .envには正しい内容が登録されている -> k8sの設定ミスではなさそう
    - ドメインはins-pre -> ins-preにつないでいるのにローカルなのはおかしい
    - ローカルで確認用に設定していたhostsファイルが原因ではないか？
- customer-apiからselect-apiに接続できない問題が発生
  - 1つ目の原因はポート
  - 上記で接続はできるようになったけど、取得に失敗する。これはおそらくaccess_keyが違うから

-

- TODOに、「ここまでやる」ラインを設ける（1Day体験からの学び）
- タイムバケットを作成する`

	ti, _ := time.Parse("2006-01-02", "2022-01-02")
	result := ExtractMemo("test-file/", 6, ti)
	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:" + result)
		t.Log("expected:" + expected)
	} else {
		t.Log("TestExtractMemo passed : \n", result)
	}
}

func TestGetCurrentQuarterFirst(t *testing.T) {
	expected := "1~3月"
	ti, _ := time.Parse("2006-01-02", "2022-01-02")
	result := GetCurrentQuarter(ti.Month())

	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:" + result)
		t.Log("expected:" + expected)
	} else {
		t.Log("TestGetCurrentQuarterFirst passed : \n", result)
	}
}

func TestGetCurrentQuarterSecond(t *testing.T) {
	expected := "4~6月"
	ti, _ := time.Parse("2006-01-02", "2022-04-02")
	result := GetCurrentQuarter(ti.Month())

	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:" + result)
		t.Log("expected:" + expected)
	} else {
		t.Log("TestGetCurrentQuarterFirst passed : \n", result)
	}
}

func TestGetCurrentQuarterThird(t *testing.T) {
	expected := "7~9月"
	ti, _ := time.Parse("2006-01-02", "2022-07-02")
	result := GetCurrentQuarter(ti.Month())

	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:" + result)
		t.Log("expected:" + expected)
	} else {
		t.Log("TestGetCurrentQuarterFirst passed : \n", result)
	}
}

func TestGetCurrentQuarterFourth(t *testing.T) {
	expected := "10~12月"
	ti, _ := time.Parse("2006-01-02", "2022-10-02")
	result := GetCurrentQuarter(ti.Month())

	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:" + result)
		t.Log("expected:" + expected)
	} else {
		t.Log("TestGetCurrentQuarterFirst passed : \n", result)
	}
}

func TestGetCurrentOKR(t *testing.T) {
	expected := `### 1~3月

- Objective1: 目標1だよ
  - Key Result1: 成果指標1
    - Score:
  - Key Result2: 成果指標2
    - Score:
- Objective2: 目標2だよ
  - Key Result1: 成果指標1
    - Score:
  - Key Result2: 成果指標2
    - Score:
- Objective3: 目標3だよ
  - Key Result1: 成果指標1
    - Score:
  - Key Result2: 成果指標2
    - Score:`

	ti, _ := time.Parse("2006-01-02", "2022-02-02")
	result := GetCurrentOKR("test-file/", strconv.Itoa(ti.Year()), GetCurrentQuarter(ti.Month()))

	if result["GoalsAndResults"] != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:" + result["GoalsAndResults"])
		t.Log("expected:" + expected)
	} else {
		t.Log("TestGetCurrentOKR passed : \n", result)
	}
}

// func TestExtractMonthlyKPT(t *testing.T) {

// 	type args struct {
// 		weeklyNotePath string
// 	}

// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{"test", args{"test-file/weekly-reviews/"}},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ExtractMonthlyKPT(tt.args.weeklyNotePath)
// 		})
// 	}
// }
func TestExtractMonthlyArticles(t *testing.T) {
	expected := `MacとLinuxそれぞれでdateコマンド使って日付計算する
https://zenn.dev/captain_blue/articles/using-the-date-on-mac-and-linux
動的コンポーネントで@clickが動かないときは@click.nativeを使う
https://zenn.dev/captain_blue/articles/nuxt-click-event-does-not-move
`

	ti, _ := time.Parse("2006-01", "2021-11")
	result := ExtractMonthlyArticles(ti)

	if result != expected {
		t.Error("抽出した文字列が想定と異なります")
		t.Log("result:" + result)
		t.Log("expected:" + expected)
	} else {
		t.Log("TestExtractMonthlyArticles passed : \n", result)
	}
}
