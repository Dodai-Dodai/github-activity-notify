package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

func main() {
	/*err := godotenv.Load("url.env")
	if err != nil {
		log.Fatal("Errload .env")
	}
	url := os.Getenv("URL")
	*/

	url := os.Getenv("URL")
	if url == " " {
		log.Fatal("Errload env:URL")
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var contributions Contributions
	json.Unmarshal(body, &contributions)

	var dataSlice []ContributionData

	for _, weeklyContributions := range contributions.Contributions {
		for _, contribution := range weeklyContributions {
			dataSlice = append(dataSlice, ContributionData{
				Date:              contribution.Date,
				ContributionCount: contribution.ContributionCount,
			})
		}
	}

	//fmt.Println(dataSlice)

	fmt.Print("昨日のコントリビューション数:")
	fmt.Println(dataSlice[len(dataSlice)-2].ContributionCount)

	var continueDays = 0
	for i := 0; i < len(dataSlice)-1; i++ {
		//昨日まで何日続いているか
		if dataSlice[i].ContributionCount == 0 {
			continueDays = 0
		} else {
			continueDays++
		}
	}

	fmt.Print("昨日までの連続コントリビューション日数:")
	fmt.Println(continueDays)

	sendLine(dataSlice[len(dataSlice)-2].ContributionCount, continueDays)
}

func sendLine(yesterday int, continueDays int) {
	/*err := godotenv.Load("line-notify.env")
	if err != nil {
		log.Fatal("Errload .env")
	}
	token := os.Getenv("TOKEN")*/

	token := os.Getenv("TOKEN")
	if token == " " {
		log.Fatal("Errload env:TOKEN")
	}

	lineURL := "https://notify-api.line.me/api/notify"

	u, err := url.ParseRequestURI(lineURL)
	if err != nil {
		log.Fatal(err)
	}

	// 昨日までのコントリビューション数と連続コントリビューション日数を送信
	text := "昨日のコントリビューション数:" + strconv.Itoa(yesterday) + "\n" + "昨日までの連続コントリビューション日数:" + strconv.Itoa(continueDays)
	// メッセージを送信
	message := url.Values{"message": {text}}
	r, _ := http.NewRequest("POST", u.String(), strings.NewReader(message.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
