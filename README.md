# genote

## はじめに

`genote`は、日次、週次、月次、年次のレビューノートを生成するためのツールです。
Zettelkasten、リサーチログ、アジャイルスタートのノートも生成できます。
コマンドラインツールとして動作します。

`genote`を使用すると、定期的なレビューノートの作成を自動化できます。
テンプレートを使用することで、ノートの形式を統一できます。
過去のノートを参照することで、振り返りを容易にできます。

## インストール

Goがインストールされていることを前提とします。

```bash
go install github.com/your-username/genote
```

`genote`の実行には、以下の環境変数が必要です。

*   `DAILY_NOTE_PATH`: 日次ノートの保存先ディレクトリ
*   `WEEKLY_NOTE_PATH`: 週次ノートの保存先ディレクトリ
*   `MONTHLY_NOTE_PATH`: 月次ノートの保存先ディレクトリ
*   `YEARLY_NOTE_PATH`: 年次ノートの保存先ディレクトリ
*   `NOTE_PATH`: Zettelkastenノートの保存先ディレクトリ
*   `DAILY_NOTE_TEMPLATE_PATH`: 日次ノートのテンプレートファイルパス
*   `WEEKLY_NOTE_TEMPLATE_PATH`: 週次ノートのテンプレートファイルパス
*   `MONTHLY_NOTE_TEMPLATE_PATH`: 月次ノートのテンプレートファイルパス
*   `ZETTELKASTEN_NOTE_TEMPLATE_PATH`: Zettelkastenノートのテンプレートファイルパス
*   `RESEARCH_LOG_NOTE_TEMPLATE_PATH`: リサーチログノートのテンプレートファイルパス
*   `AGILE_START_NOTE_TEMPLATE_PATH`: アジャイルスタートノートのテンプレートファイルパス

## 使い方

`genote`は、以下のオプションを受け付けます。

*   `-t`: テンプレートを指定します (必須)。
    *   `daily`: 日次ノート
    *   `weekly`: 週次ノート
    *   `monthly`: 月次ノート
    *   `zettelkasten`: Zettelkastenノート
    *   `research`: リサーチログノート
    *   `agile-start`: アジャイルスタートノート
    *   `memo`: メモ
*   `-d`: 日付を指定します (省略可能)。
    *   日付の形式は`YYYY-MM-DD`です。
    *   省略した場合、現在の日付が使用されます。
*   `-bd`: メモ出力でさかのぼる日数を指定します (省略可能)。
    *   デフォルト値は6です。

各テンプレートの使用例を以下に示します。

```bash
# 日次ノートを作成する
genote -t daily

# 2023年12月31日の週次ノートを作成する
genote -t weekly -d 2023-12-31

# メモを出力する (過去7日分)
genote -t memo -bd 7
```

## ライセンス

`genote`は、MITライセンスの下で公開されています。

## 開発者

Hayato Aoki

## 連絡先

`genote`に関する質問や要望は、[GitHub Issues](https://github.com/your-username/genote/issues)までお寄せください。
