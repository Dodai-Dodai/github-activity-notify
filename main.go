package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Contribution struct {
	Date  time.Time
	Count int
}

func main() {
	// github.com/Dodai-Dodaiをhttp.Getして、HTMLを取得
	// そのHTMLをパースして、Contributionのスライスを取得
	// そのスライスを出力

	// 1. github.com/Dodai-Dodaiをhttp.Getして、HTMLを取得
	url := "https://github.com/Dodai-Dodai"

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	page := string(byteArray)

	//fmt.Println(page)

	// 2. pageをパースして、Contributionのスライスを取得
	// この中にある、日付とコントリビューション数を取得する、以下は例
	/*
				<tool-tip id="tooltip-25b93dfc-4bb5-4d42-9cc1-78da76f9dd1a" for="contribution-day-component-2-49" popover="manual" data-direction="n" data-type="label" data-view-component="true" class="sr-only position-absolute">No contributions on February 27th.</tool-tip>

		                <td tabindex="0" data-ix="50" aria-selected="false" aria-describedby="contribution-graph-legend-level-1" style="width: 10px" data-date="2024-03-05" id="contribution-day-component-2-50" data-level="1" role="gridcell" data-view-component="true" class="ContributionCalendar-day"></td>
	*/

	// 3 Contributions on からコントリビューション数を取得する　このとき、コミットがない場合、数字のところにNoが入っているので、その場合は0Noを0に変換する
	ContributionPattern := `data-level="(\d)"`
	// 2024-03-05 という文字列を取得するための正規表現
	datePattern := `data-date="(\d{4}-\d{2}-\d{2})"`

	// 正規表現をコンパイル
	dateRegex := regexp.MustCompile(datePattern)
	ContributionRegex := regexp.MustCompile(ContributionPattern)

	var dates []string
	var counts []int

	dateMatches := dateRegex.FindAllStringSubmatch(page, -1)
	contributeMatches := ContributionRegex.FindAllStringSubmatch(page, -1)

	for _, dateMatch := range dateMatches {
		dates = append(dates, dateMatch[1])
	}

	fmt.Println("dates")
	for _, contributeMatch := range contributeMatches {

		counts = append(counts, int(contributeMatch[1][0]-'0'))

	}

	fmt.Println(dates)
	fmt.Println(counts)

	//Contributionを日付順に並び替える
	for i := 0; i < len(dates); i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i] > dates[j] {
				dates[i], dates[j] = dates[j], dates[i]
				counts[i], counts[j] = counts[j], counts[i]
			}
		}
	}

	for i := 0; i < len(dates); i++ {
		fmt.Printf("日付: %s, コントリビューション数: %d\n", dates[i], counts[i])
	}
	fmt.Print("昨日のコントリビューション: ")
	fmt.Println(counts[len(counts)-1])
	fmt.Print("昨日の日付: ")
	fmt.Println(dates[len(dates)-1])

	// 昨日の時点で、何日間コントリビューションが続いているかを出力
	continueDays := 0
	for i := 0; i < len(dates)-1; i++ {

		if counts[i] > 0 {
			continueDays++
		} else {
			continueDays = 0
		}
	}

	fmt.Printf("昨日の時点で、コントリビューションが続いている日数: %d\n", continueDays)

}
