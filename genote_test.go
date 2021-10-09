package main

import (
	"reflect"
	"testing"
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
	if reflect.DeepEqual(result, expected) {
		t.Error("抽出した文字列が想定と異なります")
	} else {
		t.Log("TestExtractWeeklyFDL passed : \n", result)
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
