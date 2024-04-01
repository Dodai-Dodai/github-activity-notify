package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	fmt.Print("今日のコントリビューション数:")
	fmt.Println(dataSlice[len(dataSlice)-1].ContributionCount)

	SendLine(dataSlice[len(dataSlice)-2].ContributionCount, continueDays, dataSlice[len(dataSlice)-1].ContributionCount)
}
