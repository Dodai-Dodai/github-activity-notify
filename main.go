package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// 環境変数のグローバル変数宣言
var (
	LINE_TOKEN   string
	GITHUB_TOKEN string
	GITHUB_USER  string
	URL          = "https://api.github.com/graphql"
	QUERY        string
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

func init() {
	// .envファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 環境変数の取得
	LINE_TOKEN = os.Getenv("LINE_TOKEN")
	GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
	GITHUB_USER = os.Getenv("GITHUB_USER")

	// GraphQLのクエリ
	QUERY = fmt.Sprintf(`
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
    }`, GITHUB_USER)
}

func main() {

	requestBody, err := json.Marshal(map[string]string{"query": QUERY})
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", URL, strings.NewReader(string(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", GITHUB_TOKEN))

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
