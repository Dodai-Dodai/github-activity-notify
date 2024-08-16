package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Define a struct to match the JSON structure
type Contributions struct {
	Contributions [][]struct {
		Date              string `json:"date"`
		ContributionCount int    `json:"contributionCount"`
	} `json:"contributions"`
}

// Define a struct for storing the data you need
type ContributionData struct {
	Date              string
	ContributionCount int
}

type GithubContribution struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []struct {
						ContributionDays []struct {
							Color             string `json:"color"`
							ContributionCount int    `json:"contributionCount"`
							Date              string `json:"date"`
							Weekday           int    `json:"weekday"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
}

var LINE_TOKEN = os.Getenv("LINE_TOKEN")
var TOKEN = os.Getenv("GITHUB_TOKEN")
var USER = os.Getenv("GITHUB_USER")
var URL = "https://api.github.com/graphql"
var QUERY = fmt.Sprintf(`
{
  user(login: "%s") {
    contributionsCollection {
	  contributionCalendar {
		weeks {
		  contributionDays {
			color
			contributionCount
			date
			weekday
		}
	  }
	}
  }
  }
}
`, USER)

func main() {
	requestBody, err := json.Marshal(map[string]string{"query": QUERY})
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", URL, strings.NewReader(string(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TOKEN))

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var githubContribution GithubContribution
	if err := json.NewDecoder(response.Body).Decode(&githubContribution); err != nil {
		log.Fatal(err)
	}

	yesterdayContribution := 0
	todayContribution := 0

	for _, week := range githubContribution.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			// 昨日の日付を取得
			if day.Date == time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
				fmt.Print("昨日のコントリビューション数:")
				fmt.Println(day.ContributionCount)
				yesterdayContribution = day.ContributionCount
			}
			if day.Date == time.Now().Format("2006-01-02") {
				fmt.Print("今日のコントリビューション数:")
				fmt.Println(day.ContributionCount)
				todayContribution = day.ContributionCount
			}
		}
	}

	fmt.Print("昨日までの連続コントリビューション日数:")
	var continueDays = 0
	for _, week := range githubContribution.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			if day.ContributionCount == 0 && day.Date != time.Now().Format("2006-01-02") {
				continueDays = 0
			} else if day.Date != time.Now().Format("2006-01-02") {
				continueDays++
			}
		}
	}
	fmt.Println(continueDays)

	SendLine(yesterdayContribution, continueDays, todayContribution)
}
