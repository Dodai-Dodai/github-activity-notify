package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// ContributionData 日付とコントリビューション数のペアを保持する構造体
type ContributionData struct {
	Date          time.Time
	Contributions int
}

func main() {
	url := "https://github.com/Dodai-Dodai" // GitHubのプロフィールページのURLを指定
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	html := string(body)

	var contributionsData []ContributionData

	// data-date属性と対応するコントリビューション数を抽出する正規表現
	// <td>の中の日付と、contributionsの直前の数字を取得
	re := regexp.MustCompile(`<td.+?data-date="(\d{4}-\d{2}-\d{2})".+?</td>\s*<tool-tip.+?>(\d+|No) contributions`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		dateStr := match[1]
		contributionsStr := match[2]

		date, _ := time.Parse("2006-01-02", dateStr)
		var contributions int
		if contributionsStr == "No" {
			contributions = 0
		} else {
			contributions, _ = strconv.Atoi(contributionsStr)
		}

		// スライスにコントリビューション数と日付を保存する
		contributionsData = append(contributionsData, ContributionData{
			Date:          date,
			Contributions: contributions,
		})
	}

	// 日付に基づいてソート
	sort.Slice(contributionsData, func(i, j int) bool {
		return contributionsData[i].Date.Before(contributionsData[j].Date)
	})

	// ソートされたデータの出力
	for _, data := range contributionsData {
		fmt.Printf("%s, %d\n", data.Date.Format("2006-01-02"), data.Contributions)
	}
}
